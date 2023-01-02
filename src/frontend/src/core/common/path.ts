class Path {
  private path: string[];

  constructor(base = "/") {
    this.path = ["/"];
    if (base !== "") {
      this.append(base);
    }
  }

  public append(path: string) {
    if (path == "" || path == "/") {
      return;
    }
    this.path.push(path);
  }

  public toString(): string {
    return this.path.join("/");
  }
}

export { Path };
