import { ReactNode } from "react";
import clsx from "clsx";

/** Card is the standard surface for grouped content. */
export function Card({
  className,
  children,
}: {
  className?: string;
  children: ReactNode;
}) {
  return <div className={clsx("card p-5", className)}>{children}</div>;
}
