import React from "react";
import { Link } from "react-router-dom";

import { getTitle } from "../../core/common/document";

export default function Error404() {
  getTitle("404 Not Found");

  return (
    <div>
      <h1 className="text-4xl">404 Not found</h1>
      <Link to="/">Back Home</Link>
    </div>
  );
}