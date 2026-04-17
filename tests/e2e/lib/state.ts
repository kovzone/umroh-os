import * as fs from "fs";
import * as path from "path";

const STATE_FILE = path.join(__dirname, "..", ".state.json");

interface State {
  // IAM tokens (populated by 02-iam-auth tests)
  adminToken: string;
  adminUserId: string;
}

const defaults: State = {
  adminToken: "",
  adminUserId: "",
};

export function getState(): State {
  try {
    const raw = fs.readFileSync(STATE_FILE, "utf-8");
    return { ...defaults, ...JSON.parse(raw) };
  } catch {
    return { ...defaults };
  }
}

export function setState(updates: Partial<State>): void {
  const current = getState();
  const merged = { ...current, ...updates };
  fs.writeFileSync(STATE_FILE, JSON.stringify(merged, null, 2));
}

export function clearState(): void {
  try {
    fs.unlinkSync(STATE_FILE);
  } catch {
    // ignore
  }
}
