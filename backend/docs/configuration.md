# Configuration & Secrets Management

## Overview

This project loads all configuration from environment variables and YAML files. **All secrets (API keys, passwords, tokens, etc.) must be loaded from environment variables or a secrets manager.** No default secrets should ever be present in `config.yaml` or any other config file.

## Configuration Loading

- The main config file is `configs/config.yaml`.
- All secret values in YAML must use the `${ENV_VAR}` syntax, e.g.:
  ```yaml
  api_key: "${MEXC_API_KEY}"
  ```
- At runtime, the config loader will substitute these with the values from the environment.
- If a required secret is missing from the environment, the application will fail to start and log a fatal error.
- Local development: Copy `.env.example` to `.env` and fill in required secrets. Use a tool like `direnv` or `dotenv` to load them.
- CI/Production: Inject secrets using your CI/CD pipeline or a secrets manager (e.g., AWS Secrets Manager, Vault).

## Example: Required Environment Variables

See `.env.example` for a list of all required and optional environment variables.

## Security Best Practices

- **Never** commit real secrets to source control.
- Always use environment variables or a secrets manager for sensitive data.
- Review `config.yaml` and `.env.example` before every commit to ensure no secrets are leaked.
- Rotate secrets regularly and remove unused keys.

## Troubleshooting

- If the app fails to start with a missing secret error, check your environment variables.
- Use `echo $ENV_VAR_NAME` to verify secrets are loaded in your shell.

---

_Last updated: 2025-04-15_
