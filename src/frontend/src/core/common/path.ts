class Path {
  private path: string[] = [];

  constructor(base: string | string[]) {
    this.reset(base);
  }

  public reset(path: string | string[]): Path {
    this.path = [];
    if (path.length != 0) {
      if (typeof path === "string") {
        path.split("/").forEach((value) => {
          this.append(value);
        });
      } else {
        path.forEach((value) => {
          this.append(value);
        });
      }
    }
    return this;
  }

  public append(path: string): Path {
    if (path == "" || path == "/") {
      return this;
    }
    this.path.push(path);
    return this;
  }

  public toString(): string {
    return this.path.join("/");
  }

  public forEach(
    callbackfn: (value: string, index: number, array: string[]) => void
  ): void {
    return this.path.forEach(callbackfn);
  }
}

export { Path };
