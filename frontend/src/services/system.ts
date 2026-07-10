// Frontend service wrapper around the Go SystemHandler binding.
//
// Pattern: every module gets a thin service file like this that re-exports the
// Wails-generated bindings under clean names and owns the binding import paths,
// so pages never import from ../../wailsjs directly.
import { Health } from "../../wailsjs/go/app/SystemHandler";
import { app } from "../../wailsjs/go/models";

export type SystemInfo = app.SystemInfo;

/** getHealth round-trips through Go to the database (status + shop name). */
export function getHealth(): Promise<SystemInfo> {
  return Health();
}
