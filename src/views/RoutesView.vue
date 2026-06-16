<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { z } from 'zod'
import api from '@/services/api'
import Modal from '@/components/common/Modal.vue'
import FormField from '@/components/common/FormField.vue'
import TerminalModal from '@/components/common/TerminalModal.vue'

interface ProxyRoute {
  id: number
  domain: string
  schema_type: 'static' | 'dynamic'
  target_url: string
  dynamic_resolve_url: string
  use_basic_auth: boolean
  basic_auth_username?: string
  basic_auth_password?: string
  use_validation_middleware: boolean
  validation_middleware_url?: string
  ssl_active: boolean
}

const routes = ref<ProxyRoute[]>([])
const showModal = ref(false)
const editingRoute = ref<ProxyRoute | null>(null)
const loading = ref(true)

// Terminal monitor
const showTerminal = ref(false)
const terminalRoute = ref<{ id: number; domain: string } | null>(null)

const openTerminal = (route: ProxyRoute) => {
  terminalRoute.value = { id: route.id, domain: route.domain }
  showTerminal.value = true
}

// Form Fields
const domain = ref('')
const schemaType = ref<'static' | 'dynamic'>('static')
const targetUrl = ref('')
const dynamicResolveUrl = ref('')
const useBasicAuth = ref(false)
const basicAuthUsername = ref('')
const basicAuthPassword = ref('')
const useValidationMiddleware = ref(false)
const validationMiddlewareUrl = ref('')
const sslActive = ref(false)

const errors = ref<Record<string, string>>({})

// Zod Schema
const routeSchema = z.object({
  domain: z.string().min(1, 'Domain is required').regex(/^[a-zA-Z0-9.-]+(:\d+)?(\/.*)?$/, 'Invalid domain format'),
  schema_type: z.enum(['static', 'dynamic']),
  target_url: z.string().optional().refine(val => schemaType.value === 'dynamic' || (val && val.length > 0), {
    message: 'Target URL is required for static routing'
  }),
  dynamic_resolve_url: z.string().optional().refine(val => schemaType.value === 'static' || (val && val.startsWith('http')), {
    message: 'A valid resolve URL (starting with http) is required for dynamic routing'
  }),
  basic_auth_username: z.string().optional(),
  basic_auth_password: z.string().optional(),
  validation_middleware_url: z.string().optional().refine(val => !useValidationMiddleware.value || (val && val.startsWith('http')), {
    message: 'A valid validation middleware URL (starting with http) is required'
  })
})

const fieldSchemas: Record<string, z.ZodSchema> = {
  domain: z.string().min(1, 'Domain is required').regex(/^[a-zA-Z0-9.-]+(:\d+)?(\/.*)?$/, 'Invalid domain format'),
  target_url: z.string().min(1, 'Target URL is required'),
  dynamic_resolve_url: z.string().min(1, 'Resolve URL is required').startsWith('http', 'Must start with http:// or https://'),
  validation_middleware_url: z.string().startsWith('http', 'Must start with http:// or https://'),
}

const validateField = (field: string, value: string) => {
  if (field === 'target_url' && schemaType.value === 'dynamic') {
    delete errors.value[field]
    return
  }
  if (field === 'dynamic_resolve_url' && schemaType.value === 'static') {
    delete errors.value[field]
    return
  }
  if (field === 'validation_middleware_url' && !useValidationMiddleware.value) {
    delete errors.value[field]
    return
  }

  const schema = fieldSchemas[field]
  if (!schema) return

  const result = schema.safeParse(value.trim() || undefined)
  if (!result.success) {
    const fieldErrors = result.error.flatten().fieldErrors as Record<string, string[] | undefined>
    const keys = Object.keys(fieldErrors)
    if (keys.length > 0) {
      errors.value[field] = fieldErrors[keys[0]!]?.[0] || ''
    }
  } else {
    delete errors.value[field]
  }
}

const validateForm = () => {
  errors.value = {}
  const data = {
    domain: domain.value,
    schema_type: schemaType.value,
    target_url: targetUrl.value,
    dynamic_resolve_url: dynamicResolveUrl.value,
    basic_auth_username: basicAuthUsername.value,
    basic_auth_password: basicAuthPassword.value,
    validation_middleware_url: validationMiddlewareUrl.value
  }

  const result = routeSchema.safeParse(data)
  if (!result.success) {
    const fieldErrors = result.error.flatten().fieldErrors as Record<string, string[] | undefined>
    for (const key in fieldErrors) {
      errors.value[key] = fieldErrors[key]?.[0] || ''
    }
    return false
  }
  return true
}

const fetchRoutes = async () => {
  try {
    const res = await api.get('/api/routes')
    routes.value = res.data
  } catch (err) {
    console.error('Failed to fetch routes:', err)
  } finally {
    loading.value = false
  }
}

const openAddModal = () => {
  editingRoute.value = null
  domain.value = ''
  schemaType.value = 'static'
  targetUrl.value = ''
  dynamicResolveUrl.value = ''
  useBasicAuth.value = false
  basicAuthUsername.value = ''
  basicAuthPassword.value = ''
  useValidationMiddleware.value = false
  validationMiddlewareUrl.value = ''
  sslActive.value = false
  errors.value = {}
  showModal.value = true
}

const openEditModal = (route: ProxyRoute) => {
  editingRoute.value = route
  domain.value = route.domain
  schemaType.value = route.schema_type
  targetUrl.value = route.target_url
  dynamicResolveUrl.value = route.dynamic_resolve_url
  useBasicAuth.value = route.use_basic_auth
  basicAuthUsername.value = route.basic_auth_username || ''
  basicAuthPassword.value = route.basic_auth_password || ''
  useValidationMiddleware.value = route.use_validation_middleware
  validationMiddlewareUrl.value = route.validation_middleware_url || ''
  sslActive.value = route.ssl_active
  errors.value = {}
  showModal.value = true
}

const handleSave = async () => {
  if (!validateForm()) return

  const payload = {
    domain: domain.value,
    schema_type: schemaType.value,
    target_url: targetUrl.value,
    dynamic_resolve_url: dynamicResolveUrl.value,
    use_basic_auth: useBasicAuth.value,
    basic_auth_username: useBasicAuth.value ? basicAuthUsername.value : '',
    basic_auth_password: useBasicAuth.value ? basicAuthPassword.value : '',
    use_validation_middleware: useValidationMiddleware.value,
    validation_middleware_url: useValidationMiddleware.value ? validationMiddlewareUrl.value : '',
    ssl_active: sslActive.value
  }

  try {
    if (editingRoute.value) {
      await api.put(`/api/routes/${editingRoute.value.id}`, payload)
    } else {
      await api.post('/api/routes', payload)
    }
    showModal.value = false
    fetchRoutes()
  } catch (err) {
    console.error('Failed to save route:', err)
  }
}

const handleDelete = async (id: number) => {
  if (!confirm('Are you sure you want to delete this route?')) return
  try {
    await api.delete(`/api/routes/${id}`)
    fetchRoutes()
  } catch (err) {
    console.error('Failed to delete route:', err)
  }
}

onMounted(fetchRoutes)

// Clear conditional field errors when schema type or middleware toggle changes
watch(schemaType, () => {
  delete errors.value['target_url']
  delete errors.value['dynamic_resolve_url']
  if (schemaType.value === 'static') {
    validateField('target_url', targetUrl.value)
  } else {
    validateField('dynamic_resolve_url', dynamicResolveUrl.value)
  }
})

watch(useValidationMiddleware, (val) => {
  if (!val) {
    delete errors.value['validation_middleware_url']
  } else {
    validateField('validation_middleware_url', validationMiddlewareUrl.value)
  }
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <span class="text-caption font-jetbrains-mono tracking-caption text-blue-cornflower uppercase font-medium">ROUTING</span>
        <h2 class="text-heading-sm font-semibold text-snow mt-1 tracking-tight">Proxy Routes</h2>
        <p class="text-body-sm text-ash mt-1 max-w-lg">Configure inbound domains and path prefixes to proxy traffic to target destinations.</p>
      </div>
      <button
        @click="openAddModal"
        class="px-4 py-2.5 bg-snow text-page-ink rounded-lg text-[13px] font-semibold hover:bg-ash transition-colors cursor-pointer leading-none shrink-0"
      >
        Add Proxy Route
      </button>
    </div>

    <!-- Routes List -->
    <div v-if="loading" class="text-center text-ash py-8">
      Loading proxy routes...
    </div>

    <div v-else-if="routes.length === 0" class="border border-dashed border-steel-border rounded-lg p-12 text-center text-ash">
      No routes configured yet. Click "Add Proxy Route" to get started.
    </div>

    <div v-else class="overflow-x-auto bg-card-carbon border border-steel-border rounded-lg">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="border-b border-steel-border bg-deep-coal text-caption font-jetbrains-mono tracking-caption text-ash uppercase">
            <th class="px-6 py-3 font-medium">Domain / Prefix</th>
            <th class="px-6 py-3 font-medium">Type</th>
            <th class="px-6 py-3 font-medium">Destination</th>
            <th class="px-6 py-3 font-medium">Auto SSL</th>
            <th class="px-6 py-3 font-medium text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-steel-border text-sm text-snow">
          <tr v-for="route in routes" :key="route.id" class="hover:bg-deep-coal/50 transition-colors">
            <td
              class="px-6 py-4 font-semibold font-jetbrains-mono text-blue-cornflower cursor-pointer hover:underline"
              @click="openTerminal(route)"
              title="Click to monitor live traffic"
            >{{ route.domain }}</td>
            <td class="px-6 py-4">
              <span
                class="px-2 py-0.5 rounded-[4px] text-[10px] font-jetbrains-mono font-medium uppercase tracking-wider"
                :class="route.schema_type === 'dynamic' ? 'bg-purple-950/40 border border-purple-800/40 text-purple-300' : 'bg-green-950/40 border border-green-800/40 text-green-300'"
              >
                {{ route.schema_type }}
              </span>
            </td>
            <td class="px-6 py-4 font-jetbrains-mono max-w-xs truncate text-ash">{{ route.schema_type === 'dynamic' ? route.dynamic_resolve_url : route.target_url }}</td>
            <td class="px-6 py-4">
              <span
                class="px-2 py-0.5 rounded-[4px] text-[10px] font-jetbrains-mono font-medium uppercase tracking-wider"
                :class="route.ssl_active ? 'bg-blue-cornflower/10 border border-blue-cornflower/30 text-blue-cornflower' : 'bg-graphite/20 border border-graphite text-ash'"
              >
                {{ route.ssl_active ? 'Active' : 'Disabled' }}
              </span>
            </td>
            <td class="px-6 py-4 text-right space-x-2">
              <button 
                @click="openEditModal(route)" 
                class="px-2.5 py-1 text-xs border border-graphite rounded-lg text-ash hover:text-snow hover:border-steel-border transition-colors cursor-pointer"
              >
                Edit
              </button>
              <button 
                @click="handleDelete(route.id)" 
                class="px-2.5 py-1 text-xs border border-red-900/40 bg-red-950/10 rounded-lg text-red-400 hover:text-red-300 hover:bg-red-950/20 transition-colors cursor-pointer"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Form -->
    <Modal :show="showModal" :title="editingRoute ? 'Edit Proxy Route' : 'Add Proxy Route'" @close="showModal = false">
      <div class="space-y-4">
        <FormField label="Incoming Domain / Path Prefix" id="domain" :error="errors.domain" required>
          <input
            v-model="domain"
            type="text"
            id="domain"
            @input="validateField('domain', ($event.target as HTMLInputElement).value)"
            class="w-full bg-deep-coal border border-graphite rounded-lg px-4 py-2 text-snow focus:outline-none focus:border-blue-cornflower transition-colors"
            placeholder="example.com/api"
          />
        </FormField>

        <FormField label="Routing Scheme" id="schemaType">
          <select
            v-model="schemaType"
            id="schemaType"
            class="w-full bg-deep-coal border border-graphite rounded-lg px-4 py-2 text-snow focus:outline-none focus:border-blue-cornflower transition-colors"
          >
            <option value="static">Static (Direct URL)</option>
            <option value="dynamic">Dynamic (External Resolver)</option>
          </select>
        </FormField>

        <!-- Static Fields -->
        <div v-if="schemaType === 'static'">
          <FormField label="Target Host / Destination URL" id="targetUrl" :error="errors.target_url" required>
            <input
              v-model="targetUrl"
              type="text"
              id="targetUrl"
              @input="validateField('target_url', ($event.target as HTMLInputElement).value)"
              class="w-full bg-deep-coal border border-graphite rounded-lg px-4 py-2 text-snow focus:outline-none focus:border-blue-cornflower transition-colors"
              placeholder="127.0.0.1:3000"
            />
          </FormField>
        </div>

        <!-- Dynamic Fields -->
        <div v-if="schemaType === 'dynamic'">
          <FormField label="Dynamic Resolver Endpoint" id="dynamicResolveUrl" :error="errors.dynamic_resolve_url" required>
            <input
              v-model="dynamicResolveUrl"
              type="text"
              id="dynamicResolveUrl"
              @input="validateField('dynamic_resolve_url', ($event.target as HTMLInputElement).value)"
              class="w-full bg-deep-coal border border-graphite rounded-lg px-4 py-2 text-snow focus:outline-none focus:border-blue-cornflower transition-colors"
              placeholder="https://api.mybackend.com/resolve-route"
            />
          </FormField>
        </div>

        <!-- Basic Auth Toggle -->
        <div class="border-t border-steel-border pt-4">
          <label class="flex items-center gap-2 text-sm text-snow cursor-pointer select-none">
            <input type="checkbox" v-model="useBasicAuth" class="rounded bg-deep-coal border-graphite text-blue-cornflower focus:ring-0" />
            <span>Enable Target Basic Authentication</span>
          </label>
        </div>

        <div v-if="useBasicAuth" class="grid grid-cols-2 gap-4">
          <FormField label="Username" id="authUsername">
            <input
              v-model="basicAuthUsername"
              type="text"
              id="authUsername"
              class="w-full bg-deep-coal border border-graphite rounded-lg px-4 py-2 text-snow focus:outline-none focus:border-blue-cornflower transition-colors"
            />
          </FormField>
          <FormField label="Password" id="authPassword">
            <input
              v-model="basicAuthPassword"
              type="password"
              id="authPassword"
              class="w-full bg-deep-coal border border-graphite rounded-lg px-4 py-2 text-snow focus:outline-none focus:border-blue-cornflower transition-colors"
            />
          </FormField>
        </div>

        <!-- Validation Middleware Toggle -->
        <div class="border-t border-steel-border pt-4">
          <label class="flex items-center gap-2 text-sm text-snow cursor-pointer select-none">
            <input type="checkbox" v-model="useValidationMiddleware" class="rounded bg-deep-coal border-graphite text-blue-cornflower focus:ring-0" />
            <span>Enable Validation Middleware</span>
          </label>
        </div>

        <div v-if="useValidationMiddleware">
          <FormField label="Validation Endpoint URL" id="validationUrl" :error="errors.validation_middleware_url" required>
            <input
              v-model="validationMiddlewareUrl"
              type="text"
              id="validationUrl"
              @input="validateField('validation_middleware_url', ($event.target as HTMLInputElement).value)"
              class="w-full bg-deep-coal border border-graphite rounded-lg px-4 py-2 text-snow focus:outline-none focus:border-blue-cornflower transition-colors"
              placeholder="https://auth.mybackend.com/validate"
            />
          </FormField>
        </div>

        <!-- Auto SSL Toggle -->
        <div class="border-t border-steel-border pt-4">
          <label class="flex items-center gap-2 text-sm text-snow cursor-pointer select-none">
            <input type="checkbox" v-model="sslActive" class="rounded bg-deep-coal border-graphite text-blue-cornflower focus:ring-0" />
            <span>Enable Auto Let's Encrypt SSL (Requires domain DNS mapped to server)</span>
          </label>
        </div>
      </div>

      <template #footer>
        <button
          @click="showModal = false"
          class="px-4 py-2 border border-graphite rounded-lg text-snow text-sm font-medium hover:bg-card-carbon transition-colors cursor-pointer"
        >
          Cancel
        </button>
        <button
          @click="handleSave"
          class="px-4 py-2 bg-snow text-page-ink font-medium rounded-lg text-sm hover:bg-ash transition-colors cursor-pointer"
        >
          Save Route
        </button>
      </template>
    </Modal>

    <!-- Terminal Monitor Modal -->
    <TerminalModal
      v-if="terminalRoute"
      :show="showTerminal"
      :route-id="terminalRoute.id"
      :domain="terminalRoute.domain"
      @close="showTerminal = false"
    />
  </div>
</template>
