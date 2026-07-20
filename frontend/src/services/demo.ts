import { LoadDemoData } from "../../wailsjs/go/app/DemoHandler";
import { service } from "../../wailsjs/go/models";

export type DemoSummary = service.DemoSummary;

export const demoApi = {
  load: () => LoadDemoData(),
};
