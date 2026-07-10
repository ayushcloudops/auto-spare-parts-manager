import { useEffect, useState } from "react";
import { getHealth, SystemInfo } from "../services/system";

interface HealthState {
  info: SystemInfo | null;
  error: boolean;
}

/** useHealth fetches backend health once on mount. */
export function useHealth(): HealthState {
  const [state, setState] = useState<HealthState>({ info: null, error: false });

  useEffect(() => {
    let active = true;
    getHealth()
      .then((info) => active && setState({ info, error: false }))
      .catch(() => active && setState({ info: null, error: true }));
    return () => {
      active = false;
    };
  }, []);

  return state;
}
