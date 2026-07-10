import { Construction } from "lucide-react";

/**
 * ComingSoon is a placeholder body for modules not yet implemented. Each module
 * step replaces its page's use of this with the real screen.
 */
export default function ComingSoon({ title }: { title: string }) {
  return (
    <div className="card flex flex-col items-center justify-center gap-3 p-12 text-center">
      <div className="flex h-12 w-12 items-center justify-center rounded-full bg-slate-100 text-slate-400 dark:bg-slate-800">
        <Construction size={22} />
      </div>
      <div>
        <p className="font-medium text-slate-700 dark:text-slate-200">
          {title} module
        </p>
        <p className="text-sm text-slate-400">Coming up in a later build step.</p>
      </div>
    </div>
  );
}
