import { PreviewReceipt, PrintReceipt } from "../../wailsjs/go/app/PrintHandler";

export const printApi = {
  preview: (invoiceId: number) => PreviewReceipt(invoiceId),
  print: (invoiceId: number) => PrintReceipt(invoiceId),
};
