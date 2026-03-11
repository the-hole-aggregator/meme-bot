## 🔧 Git Workflow

We follow the [Gitflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow) branching strategy.

---

### 🧞 Commit Naming Convention

Use a simplified [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0) format:

```
MEME<TASK-ID> - <type>(<scope>): <short summary>
```

#### Example:
```
MEME-12345: refactor(api): change endpoint name
```

- **TASK-ID** — automatically added via [git hook](/.githooks/prepare-commit-msg)
- **type**:
  - `feat` — new feature
  - `fix` — bug fix
  - `refactor` — code refactoring
  - `docs` — documentation changes
  - `test` — tests added or updated
  - `perf` — performance optimization
  - `style` — formatting/stylistic changes
  - `build` — build-related changes
  - `ci` — CI configuration changes
  - `chore` — miscellaneous, non-code changes

- **scope** — feature/module/class/function where changes were made
- **short summary** — must use the **infinitive form**, e.g., `add filter`, not `filter added`

> 💡 **One commit = one atomic change.**  
> Avoid bundling unrelated changes into a single commit.

---

### 🌿 Branches

- `main` — latest production-ready version (what users currently have)
- `dev` — all tested features merged here; always stable
- `feature/MEME-<task-id>` — new feature development
- `release/<date>` — created from `dev` for release preparation and QA

---

### ✅ Code Review Rules

- Use **Merge** when merging PRs
- If changes are needed after PR approval:
  1. Original: `0001-feature-some` → `dev` ✅
  2. Fix: `0001-fix-some` → `0001-feature-some` 🔁

---

### 🔁 Keeping Feature Branch Up-to-Date

Always rebase your branch on top of `dev` before merging or creating a PR:

### 🚀 After Release

1. Merge release into `main`:

2. Delete the release branch:

3. Merge `main` into `dev` to sync history:

---