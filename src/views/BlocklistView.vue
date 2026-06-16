<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { z } from 'zod'
import api from '@/services/api'
import Modal from '@/components/common/Modal.vue'
import FormField from '@/components/common/FormField.vue'

interface BlockedIP {
  id: number
  ip_address: string
  reason: string
  created_at: string
}

const blocklist = ref<BlockedIP[]>([])
const showModal = ref(false)
const loading = ref(true)

// Form Fields
const ipAddress = ref('')
const reason = ref('')
const errors = ref<Record<string, string>>({})

// Zod Schema
const ipSchema = z.object({
  ip_address: z.string().min(1, 'IP Address is required').refine(val => {
    // Basic IPv4 regex
    const ipv4Regex = /^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$/
    // Basic IPv6 regex
    const ipv6Regex = /^([\da-fA-F]{1,4}:){7}[\da-fA-F]{1,4}$/
    return ipv4Regex.test(val) || ipv6Regex.test(val)
  }, {
    message: 'Invalid IP Address format'
  }),
  reason: z.string().min(1, 'Reason is required')
})

const validateForm = () => {
  errors.value = {}
  const result = ipSchema.safeParse({
    ip_address: ipAddress.value,
    reason: reason.value
  })

  if (!result.success) {
    const fieldErrors = result.error.flatten().fieldErrors as Record<string, string[] | undefined>
    for (const key in fieldErrors) {
      errors.value[key] = fieldErrors[key]?.[0] || ''
    }
    return false
  }
  return true
}

const fetchBlocklist = async () => {
  try {
    const res = await api.get('/api/blocklist')
    blocklist.value = res.data
  } catch (err) {
    console.error('Failed to fetch blocklist:', err)
  } finally {
    loading.value = false
  }
}

const openAddModal = () => {
  ipAddress.value = ''
  reason.value = ''
  errors.value = {}
  showModal.value = true
}

const handleSave = async () => {
  if (!validateForm()) return

  try {
    await api.post('/api/blocklist', {
      ip_address: ipAddress.value,
      reason: reason.value
    })
    showModal.value = false
    fetchBlocklist()
  } catch (err: any) {
    if (err.response?.data?.includes('UNIQUE')) {
      errors.value.ip_address = 'This IP address is already blocked'
    } else {
      console.error('Failed to block IP:', err)
    }
  }
}

const handleUnblock = async (id: number) => {
  if (!confirm('Are you sure you want to unblock this IP address?')) return
  try {
    await api.delete(`/api/blocklist/${id}`)
    fetchBlocklist()
  } catch (err) {
    console.error('Failed to unblock IP:', err)
  }
}

onMounted(fetchBlocklist)
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <span class="text-caption font-jetbrains-mono tracking-caption text-blue-cornflower uppercase font-medium">SECURITY</span>
        <h2 class="text-3xl font-semibold text-snow mt-1 tracking-tight">IP Block List</h2>
      </div>
      <button
        @click="openAddModal"
        class="px-4 py-2 bg-snow text-page-ink font-medium rounded-lg text-sm hover:bg-ash transition-colors cursor-pointer"
      >
        Block IP Address
      </button>
    </div>

    <!-- Blocklist -->
    <div v-if="loading" class="text-center text-ash py-8">
      Loading IP block list...
    </div>

    <div v-else-if="blocklist.length === 0" class="border border-dashed border-steel-border rounded-lg p-12 text-center text-ash">
      No IP addresses are currently blocked.
    </div>

    <div v-else class="overflow-x-auto bg-card-carbon border border-steel-border rounded-lg">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="border-b border-steel-border bg-deep-coal text-caption font-jetbrains-mono tracking-caption text-ash uppercase">
            <th class="px-6 py-3 font-medium">IP Address</th>
            <th class="px-6 py-3 font-medium">Blocked Reason</th>
            <th class="px-6 py-3 font-medium">Blocked At</th>
            <th class="px-6 py-3 font-medium text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-steel-border text-sm text-snow">
          <tr v-for="item in blocklist" :key="item.id" class="hover:bg-deep-coal/50 transition-colors">
            <td class="px-6 py-4 font-semibold font-jetbrains-mono text-red-400">{{ item.ip_address }}</td>
            <td class="px-6 py-4 text-ash">{{ item.reason }}</td>
            <td class="px-6 py-4 text-ash font-jetbrains-mono text-xs">
              {{ new Date(item.created_at).toLocaleString() }}
            </td>
            <td class="px-6 py-4 text-right">
              <button 
                @click="handleUnblock(item.id)" 
                class="px-2.5 py-1 text-xs border border-red-900/40 bg-red-950/10 rounded-lg text-red-400 hover:text-red-300 hover:bg-red-950/20 transition-colors cursor-pointer"
              >
                Unblock
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Form -->
    <Modal :show="showModal" title="Block IP Address" @close="showModal = false">
      <div class="space-y-4">
        <FormField label="IP Address" id="ipAddress" :error="errors.ip_address" required>
          <input
            v-model="ipAddress"
            type="text"
            id="ipAddress"
            class="w-full bg-deep-coal border border-graphite rounded-lg px-4 py-2 text-snow focus:outline-none focus:border-blue-cornflower transition-colors font-jetbrains-mono"
            placeholder="192.168.1.100"
          />
        </FormField>

        <FormField label="Reason for Blocking" id="reason" :error="errors.reason" required>
          <input
            v-model="reason"
            type="text"
            id="reason"
            class="w-full bg-deep-coal border border-graphite rounded-lg px-4 py-2 text-snow focus:outline-none focus:border-blue-cornflower transition-colors"
            placeholder="DDoS source, spam request, etc."
          />
        </FormField>
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
          Block IP
        </button>
      </template>
    </Modal>
  </div>
</template>
