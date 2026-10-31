// Pi's side of the plugin: the contract, notebook install and correction flagging that
// hooks/hooks.json gives Claude Code and Codex, plus the researcher agent,
// which a Pi package cannot carry. Skills come from package.json.
import { execFile, spawn } from "node:child_process";
import { existsSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { homedir } from "node:os";
import { join } from "node:path";
import { fileURLToPath } from "node:url";
import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";

const contract = readFileSync(new URL("../contract.md", import.meta.url), "utf8");
const flagCorrection = fileURLToPath(new URL("../hooks/flag-correction.sh", import.meta.url));
const installNotebook = fileURLToPath(new URL("../hooks/install-notebook.sh", import.meta.url));
const sessionContext = fileURLToPath(new URL("../hooks/session-context.sh", import.meta.url));

const managed = "workbench_managed: true";
const researcher = readFileSync(new URL("../agents/researcher.md", import.meta.url), "utf8")
	.replace("model: inherit\n", "")
	.replace("effort: medium", "thinking: medium")
	.replace(
		/^tools: .*$/m,
		`tools: web_search, fetch_content, read, grep, find, ls\nprompt_mode: replace\n${managed}`,
	);

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

// pi-subagents reads agents only from this folder and the project's. A
// researcher.md without the managed line is the user's own and stays.
function installResearcher() {
	const dir = join(process.env.PI_CODING_AGENT_DIR ?? join(homedir(), ".pi", "agent"), "agents");
	const path = join(dir, "researcher.md");
	const current = existsSync(path) ? readFileSync(path, "utf8") : undefined;
	if (current === researcher || (current !== undefined && !current.includes(managed))) return;
	mkdirSync(dir, { recursive: true });
	writeFileSync(path, researcher);
}

export default function workbench(pi: ExtensionAPI) {
	let context = contract;

	pi.on("session_start", (_event, ctx) => {
		installResearcher();
		// Adds the notebook push; on any failure the plain contract stays.
		execFile("sh", [sessionContext], { cwd: ctx?.cwd, timeout: 5000 }, (error, stdout) => {
			if (!error) context = stdout;
		});
		spawn("sh", [installNotebook], { stdio: "ignore" }).on("error", () => {});
	});

	pi.on("before_agent_start", (event) => {
		event.systemPromptOptions.sections.workbench_contract = context;
	});

	pi.on("input", async (event, ctx) => {
		if (event.source === "extension") return;
		await flag(event.text, ctx.sessionManager.getSessionId(), ctx.cwd);
	});
}
