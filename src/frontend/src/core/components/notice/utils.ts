import i18n from "../../i18n/i18n";

const regex = /<[^>]+>/s;
const sep = ".";
const modelKey = "{model}";
const fieldKey = "{field}";
const reasonKey = "{reason}";

interface Error {
  Error: string;
  Model: string;
  Field: string;
  Reason: string;
}

export function parseTemplate(context: string): string {
  if (!regex.test(context)) {
    return context;
  } else {
    const m = regex.exec(context);
    if (m !== null) {
      context = m[0];
      context = context.substr(1, context.length - 2);
    }
  }

  const seps = context.split(sep);
  if (seps.length == 0) {
    return context;
  }

  const keystone: string = seps[0];

  const msgTemplate: string = i18n.t("message.error." + keystone);

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

  return format(msgTemplate, err);
}

function format(temp: string, err: Error): string {
  if (temp.includes(modelKey)) {
    temp = temp.replace(modelKey, err.Model!);
  }
  if (temp.includes(fieldKey)) {
    temp = temp.replace(fieldKey, err.Field!);
  }
  if (temp.includes(reasonKey)) {
    temp = temp.replace(reasonKey, err.Reason!);
  }
  return temp;
}
