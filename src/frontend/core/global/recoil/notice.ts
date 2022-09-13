import { atom, SetterOrUpdater, useSetRecoilState } from "recoil";
import i18n from "../../i18n/i18n";

type typeNotice = "model";
type colorNotice = "primary" | "success" | "warning" | "danger";

export type noticeValue = {
  id: string;
  type: typeNotice;
  color: colorNotice;
  title: string;
  text: string;
};

export const noticeState = atom<noticeValue>({
  key: "notice",
  default: undefined,
});

export const useNotice = (set: SetterOrUpdater<noticeValue>) => {
  return new Notice(set);
};

class Notice {
  set: SetterOrUpdater<noticeValue>;

  constructor(s: SetterOrUpdater<noticeValue>) {
    this.set = s;
    return;
  }

  public primary(text: string) {
    this.emit("primary", text);
  }

  public success(text: string) {
    this.emit("success", text);
  }

  public error(text: string) {
    this.emit("danger", text);
  }

  public warning(text: string) {
    this.emit("warning", text);
  }

  public emit(color: colorNotice, text: string) {
    this.set({
      id: parser.getID(),
      type: "model",
      color: color,
      title: i18n.t("notice.title"),
      text: parser.execute(text),
    });
  }
}

interface Error {
  Error: string;
  Model: string;
  Field: string;
  Reason: string;
}

class TemplateParser {
  private static regex = /<[^>]+>/s;
  private static sep = ".";
  private static modelKey = "{model}";
  private static fieldKey = "{field}";
  private static reasonKey = "{reason}";

  private idCount = 0;

  constructor() {
    return;
  }

  public execute(context: string): string {
    if (!TemplateParser.regex.test(context)) {
      return context;
    } else {
      const m = TemplateParser.regex.exec(context);
      if (m !== null) {
        context = m[0];
        context = context.substr(1, context.length - 2);
      }
    }

    const seps = context.split(TemplateParser.sep);
    if (seps.length === 0) {
      return context;
    }

    const keystone = seps[0];
    const msgTemplate = i18n.t("message.error." + keystone);
    if (msgTemplate === null) {
      return context;
    }

    const err: Error = {
      Error: msgTemplate,
      Model: "",
      Field: "",
      Reason: "",
    };

    for (let i = 1; i < seps.length; i++) {
      if (seps[i].length == 0) continue;

      const modelPath = "message.model." + seps[i];
      const fieldPath = "message.field." + seps[i];
      const reasonPath = "message.reason." + seps[i];

      if (i18n.exists(modelPath)) {
        err.Model = i18n.t(modelPath);
      } else if (i18n.exists(fieldPath)) {
        err.Field = i18n.t(fieldPath);
      } else if (i18n.exists(reasonPath)) {
        err.Reason = i18n.t(reasonPath);
      }
    }

    return this.format(msgTemplate, err);
  }

  private format(temp: string, err: Error): string {
    if (temp.includes(TemplateParser.modelKey)) {
      temp = temp.replace(TemplateParser.modelKey, err.Model!);
    }
    if (temp.includes(TemplateParser.fieldKey)) {
      temp = temp.replace(TemplateParser.fieldKey, err.Field!);
    }
    if (temp.includes(TemplateParser.reasonKey)) {
      temp = temp.replace(TemplateParser.reasonKey, err.Reason!);
    }
    return temp;
  }

  public getID(): string {
    this.idCount++;
    return this.idCount.toString();
  }
}

const parser = new TemplateParser();
