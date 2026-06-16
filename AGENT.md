# Agent Instructions & Guidelines (AGENT.md)

Welcome! You are working on the **Dynamic Golang Proxy** project. Please read and follow the instructions below carefully.

---

## 🛠️ Custom Agent Skills & Rules
You **MUST** read and understand all custom skills and rules in the `.agents/` directory before proceeding with any development tasks:
1. **Repository Rules**:
   - [antigravity-rtk-rules.md](./.agents/rules/antigravity-rtk-rules.md) - Explains `rtk` wrapper usage to optimize token consumption.
2. **Skills in `.agents/skills/`**:
   - [caveman](./.agents/skills/caveman/SKILL.md) - Ultra-compressed communication mode.
   - [caveman-commit](./.agents/skills/caveman-commit/SKILL.md) - Conventional commits generator.
   - [caveman-compress](./.agents/skills/caveman-compress/SKILL.md) - Compressed natural language memory files.
   - [caveman-help](./.agents/skills/caveman-help/SKILL.md) - Caveman mode quick reference.
   - [caveman-review](./.agents/skills/caveman-review/SKILL.md) - Compressed review comments.
   - [compress](./.agents/skills/compress/SKILL.md) - File compressor.
   - [ui-ux-pro-max](./.agents/skills/ui-ux-pro-max/SKILL.md) - UI/UX style guide.
   - [vibe-security](./.agents/skills/vibe-security/SKILL.md) - Security vulnerability audit.
   - [vue-best-practices](./.agents/skills/vue-best-practices/SKILL.md) - Standard Vue 3 Composition API & TypeScript practices.
   - [vue-options-api-best-practices](./.agents/skills/vue-options-api-best-practices/SKILL.md) - Vue Options API reference.
   - [vue-pinia-best-practices](./.agents/skills/vue-pinia-best-practices/SKILL.md) - Pinia store management.
   - [vue-router-best-practices](./.agents/skills/vue-router-best-practices/SKILL.md) - Router configuration and lifecycle.

---

## 📋 Project Specifications
Before writing any code or modifying the application, read [SPEC.md](./SPEC.md) to understand the tech stack, UI/UX guidelines, design parameters, and features of the proxy engine.

---

## 📝 Task Workflow & TODO.md
Tasks, features, and fixes are tracked in [TODO.md](./TODO.md).

### Structure of `TODO.md`
The file consists of 3 distinct sections:
1. **Backlog**: Features to be implemented. **Do NOT use checkboxes `[ ]`** in this section; list items using standard bullet points (`- Task Name`).
2. **Error / Bug**: Active issues/bugs to be resolved. **Do NOT use checkboxes `[ ]`** in this section; list items using standard bullet points (`- Bug Name`).
3. **Finish**: Successfully completed tasks. This is the **ONLY** section that uses checkboxes, marked with a check (`- [x] Task/Bug Name`).

If the **Backlog** or **Error / Bug** sections are empty, do not leave them completely blank; instead, write `- ` (a dash followed by a space) under the header.

### Agent Workflow & Testing
- **Read First**: At the start of your execution, always inspect `TODO.md` to see what needs to be worked on.
- **Process**: Pick tasks from the **Backlog** or **Error / Bug** sections and implement them.
- **Dependencies (FE)**: If you need to add or install any Frontend (FE) dependencies, **always use `bun`** (e.g., `bun add <package-name>`). Do NOT use `npm`, `yarn`, or `pnpm` for installation.
- **Post-Prompt Verification**:
  - **If Frontend (FE) code was generated/modified**: Always test and verify changes by running:
    ```bash
    rtk yarn lint
    rtk yarn build
    rtk yarn doctor
    ```
    *(Note: `yarn doctor` runs `@rekl0w/vue-doctor` to audit Vue files)*.
    **CRITICAL REQUIREMENT**: The `yarn doctor` score must be **100% (Perfect)**. Even warnings must be fully resolved. After running `yarn doctor`, you must immediately read [vue-doctor-report.md](./vue-doctor-report.md) and resolve all reported errors and warnings until a 100% score is achieved.
  - **If Backend (BE) code was generated/modified**: Always test and verify changes by running:
    ```bash
    go build
    ```
- **Update**: Once a task is completed, verified, and passes the required tests, remove it from the **Backlog** or **Error / Bug** sections, format it as a completed checkbox (`- [x] Task/Bug Name`), and move it into the **Finish** section of `TODO.md`. If a change or feature is requested directly from chat or implemented outside of the pre-existing `TODO.md` list, make sure to add it directly to the **Finish** section of `TODO.md` as well. Keep the file organized.
- **No Mistakes**: Follow the designated tech stack, guidelines, and routing requirements exactly.
