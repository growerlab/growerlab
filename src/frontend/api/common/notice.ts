import { toaster } from 'evergreen-ui';

import i18n from '../../i18n/i18n';

const regex = /<[^>]+>/s;
const Sep = '.';
const ModelKey = '{model}';
const FieldKey = '{field}';
const ReasonKey = '{reason}';

export const Message = {
  Success: function (text: string): void {
    text = parseTemplate(text);
    return toaster.success(text);
  },
  Error: function (text: string): void {
    text = parseTemplate(text);
    return toaster.danger(text);
  },
  Warning: function (text: string): void {
    text = parseTemplate(text);
    return toaster.warning(text);
  },
};

interface Error {
  Error: string;
  Model: string;
  Field: string;
  Reason: string;
}

function parseTemplate(context: string): string {
  if (!regex.test(context)) {
    return context;
  } else {
    let m = regex.exec(context);

    if ((m = regex.exec(context)) !== null) {
      context = m[0];
      context = context.substr(1, context.length - 2);
    }
  }

  const seps = context.split(Sep);
  if (seps.length == 0) {
    return context;
  }

  const keystone: string = seps[0];

  const msgTemplate: string = i18n.t('message.error.' + keystone);

  if (msgTemplate === null) {
    return context;
  }

  const err: Error = {
    Error: msgTemplate,
    Model: '',
    Field: '',
    Reason: '',
  };

  for (let i = 1; i < seps.length; i++) {
    if (seps[i].length == 0) continue;

    const modelPath = 'message.model.' + seps[i];
    const fieldPath = 'message.field.' + seps[i];
    const reasonPath = 'message.reason.' + seps[i];

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
  if (temp.includes(ModelKey)) {
    temp = temp.replace(ModelKey, err.Model!);
  }
  if (temp.includes(FieldKey)) {
    temp = temp.replace(FieldKey, err.Field!);
  }
  if (temp.includes(ReasonKey)) {
    temp = temp.replace(ReasonKey, err.Reason!);
  }
  return temp;
}
