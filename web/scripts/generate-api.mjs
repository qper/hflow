import { existsSync } from 'node:fs';
import { execFileSync } from 'node:child_process';

const candidates = ['./internal/api/openapi.json', '../internal/api/openapi.json'];
const specPath = candidates.find((candidate) => existsSync(candidate));

if (!specPath) {
  console.error('OpenAPI spec not found. Expected one of:', candidates.join(', '));
  process.exit(1);
}

const command = process.platform === 'win32' ? 'npx.cmd' : 'npx';
execFileSync(command, ['openapi-typescript', specPath, '--output', 'src/api/generated/schema.ts'], {
  stdio: 'inherit',
});
