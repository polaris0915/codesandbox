## Git commit 提交信息的规范

**`<type>`**: 本次提交的类型，常见类型如下：

- **feat**: 添加新功能
- **fix**: 修复 Bug
- **docs**: 文档变更
- **style**: 格式变动（代码格式化、空格、注释等），不影响功能
- **refactor**: 重构代码（既不是新功能也不是修复 Bug）
- **test**: 添加或修改测试
- **chore**: 其他变更（例如构建流程、依赖更新等）
- **perf**: 性能优化

**`<scope>`**: 变更的范围（可选），如 `auth`、`login`、`backend`。

**`<subject>`**: 简短说明变更内容。

**`<body>`**: 更详细的描述（可选），可以包括：

- 修改的原因。
- 修改的方式。
- 注意事项。

**`<footer>`**: 关联信息（可选），例如：

- `Closes #123`（自动关闭 Issue #123）
- `BREAKING CHANGE:` 描述兼容性破坏的变更。

## 示例

### **提交信息示例：**

```
markdown


Copy code
feat(auth): implement JWT authentication

This update introduces a new authentication mechanism based on JWT (JSON Web Tokens) for secure user login. The changes include:

- Added `login` and `logout` API endpoints in the `auth` service.
- JWT token is now used for user sessions instead of traditional sessions.
- Updated the login route to handle the JWT token generation and validation.
- Integrated `bcrypt` for password hashing to enhance security.

Users will now authenticate via JWT, improving the scalability and performance of the authentication system.

BREAKING CHANGE: The login and session management has changed, existing authentication methods will no longer work.

Closes #45
```

------

### **解析：**

1. **<type>**: `feat`
   - 表示这次提交是 **添加新功能**。这里是新增了基于 JWT 的认证功能。
2. **<scope>**: `auth`
   - `auth` 是本次提交修改的范围。这里指的是认证系统（如登录、用户验证等）。
3. **<subject>**: `implement JWT authentication`
   - 简短描述了本次提交的变更内容，即实现了 JWT 认证。
4. **<body>**:
   - **修改的原因**：为什么做这次修改，比如为了提升认证系统的安全性和性能，换用 JWT 代替传统的会话管理。
   - **修改的方式**：具体做了什么更改，包括新增 API、集成 JWT、使用 `bcrypt` 等。
   - **注意事项**：例如用户身份验证方式变更带来的影响，解释了为什么做这次修改。
5. **<footer>**:
   - **BREAKING CHANGE**: 提示这次变更是 **破坏性变更**，意味着旧的认证方式将不再支持，需要用户适应新的 JWT 登录方式。
   - **Closes #45**: 关联 Issue，表示本次提交解决了 GitHub 中的 Issue #45，这个 Issue 会被自动关闭。

---

