#!/bin/sh
# Drives pi/workbench.ts through a stub of Pi's extension API, so it needs Node
# but no Pi install and no model.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)

if ! command -v node >/dev/null 2>&1; then
  echo "pi extension tests: skipped, node not found"
  exit 0
fi

data=$(mktemp -d)
trap 'rm -rf "$data"' EXIT

XDG_DATA_HOME="$data" PI_CODING_AGENT_DIR="$data/pi" EXTENSION="$root/pi/workbench.ts" CONTRACT="$root/contract.md" \
  node --no-warnings --experimental-strip-types --input-type=module -e '
import assert from "node:assert/strict";
import { existsSync, readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const { default: extension } = await import(process.env.EXTENSION);
const handlers = {};
extension({ on: (name, handler) => { handlers[name] = handler; } });

const prompt = { systemPromptOptions: { sections: {} } };
await handlers.before_agent_start(prompt);
assert.equal(prompt.systemPromptOptions.sections.workbench_contract,
  readFileSync(process.env.CONTRACT, "utf8"));
console.log("ok    contract section");

const ctx = { cwd: "/work", sessionManager: { getSessionId: () => "s1" } };
const pending = join(process.env.XDG_DATA_HOME, "workbench/state/pending-corrections.jsonl");

await handlers.input({ text: "add a test for this", source: "interactive" }, ctx);
assert.equal(existsSync(pending), false);
console.log("ok    plain prompt not flagged");

await handlers.input({ text: "no, that is wrong", source: "extension" }, ctx);
assert.equal(existsSync(pending), false);
console.log("ok    extension input not flagged");

await handlers.input({ text: "no, that is wrong", source: "interactive" }, ctx);
const entry = JSON.parse(readFileSync(pending, "utf8"));
assert.equal(entry.prompt, "no, that is wrong");
assert.equal(entry.session_id, "s1");
assert.equal(entry.cwd, "/work");
console.log("ok    correction flagged");

const agents = join(process.env.PI_CODING_AGENT_DIR, "agents");
const agent = join(agents, "researcher.md");

await handlers.session_start();
const written = readFileSync(agent, "utf8");
assert.match(written, /^tools: web_search, fetch_content, read, grep, find, ls$/m);
assert.match(written, /^thinking: medium$/m);
assert.match(written, /^prompt_mode: replace$/m);
assert.doesNotMatch(written, /WebSearch|^model:|^effort:/m);
assert.match(written, /Every claim cites a source/);
console.log("ok    researcher installed");

writeFileSync(agent, written.replace("Every claim", "An older claim"));
await handlers.session_start();
assert.equal(readFileSync(agent, "utf8"), written);
console.log("ok    managed copy updated");

writeFileSync(agent, "my own agent\n");
await handlers.session_start();
assert.equal(readFileSync(agent, "utf8"), "my own agent\n");
console.log("ok    personal researcher kept");
'
echo "pi extension tests: passed"
