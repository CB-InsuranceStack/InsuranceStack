import { useState, useEffect } from 'react';
import { getFlagsSnapshot, subscribeFlags } from '../features/flags';

/**
 * Reactive hook for CloudBees Feature Management flags
 * Updates automatically when flag values change in FM dashboard
 *
 * @param key - The flag key to watch
 * @returns Current boolean value of the flag
 *
 * @example
 * const isEnabled = useRoxFlag('transactionsFilters');
 */
export default function useRoxFlag(key: string): boolean {
  const [val, setVal] = useState(() => !!getFlagsSnapshot()[key]);

  useEffect(() => {
    return subscribeFlags((_reason, snap) => {
      const newVal = !!snap[key];
      // Only update if the value actually changed to prevent unnecessary re-renders
      setVal((prevVal) => (prevVal !== newVal ? newVal : prevVal));
    });
  }, [key]);

  return val;
}
