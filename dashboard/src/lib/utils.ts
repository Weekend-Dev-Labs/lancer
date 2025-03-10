import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}


export function formatText(input: string) {
  // Insert a space before all caps and trim any leading/trailing spaces
  let formatted = input.replace(/([A-Z])/g, " $1").trim();

  // Capitalize the first letter of each word
  formatted = formatted.replace(/\b\w/g, (char) => char.toUpperCase());

  return formatted;
}