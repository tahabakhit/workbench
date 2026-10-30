// Pi's side of the plugin: the contract and correction flagging that
// hooks/hooks.json gives Claude Code and Codex. Skills come from package.json.
import { spawn } from "node:child_process";
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";

const contract = readFileSync(new URL("../contract.md", import.meta.url), "utf8");
const flagCorrection = fileURLToPath(new URL("../hooks/flag-correction.sh", import.meta.url));

// A failed flag must never stop the prompt, so every outcome resolves.
function flag(prompt: string, session_id: string, cwd: string): Promise<void> {
	return new Promise((resolve) => {
		const hook = spawn("sh", [flagCorrection], { stdio: ["pipe", "ignore", "ignore"] });
		hook.on("error", () => resolve());
		hook.on("close", () => resolve());
		hook.stdin.on("error", () => {});
		hook.stdin.end(JSON.stringify({ prompt, session_id, cwd }));
	});
}

export default function workbench(pi: ExtensionAPI) {
	pi.on("before_agent_start", (event) => {
		event.systemPromptOptions.sections.workbench_contract = contract;
	});

	pi.on("input", async (event, ctx) => {
		if (event.source === "extension") return;
		await flag(event.text, ctx.sessionManager.getSessionId(), ctx.cwd);
	});
}
