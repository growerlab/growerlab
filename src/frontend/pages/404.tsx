import React from "react";
import { getTitle } from "../common/document";
import Link from "next/link";

export default function () {
  getTitle("404 Not Found");

  return (
    <div>
      <h1 className="text-4xl">404 Not found</h1>
      <Link href="/">Back Home</Link>
    </div>
  );
}
