import React from "react";
import Link from "next/link";
import { getTitle } from "../core/common/document";

export default function fn404() {
  getTitle("404 Not Found");

  return (
    <div>
      <h1 className="text-4xl">404 Not found</h1>
      <Link href="/">Back Home</Link>
    </div>
  );
}
