import React from "react";
import { useNavigate } from "react-router-dom";

import LoginForm from "../../core/components/user/LoginForm";
import i18n from "../../core/i18n/i18n";
import { Router } from "../../config/router";
import { useTitle } from "../../core/global/state";

export default function Login(props: any) {
  useTitle(i18n.t("website.login"));
  const navigate = useNavigate();

  const onSuccess = () => {
    navigate(Router.User.Index);
  };

  return (
    <>
      <div className="min-h-full flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-md w-full space-y-8">
          <div>
            <h1 className="mx-auto h-12 w-auto text-center text-3xl">Logo</h1>
            <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
              Sign in to your account
            </h2>
          </div>
          <div>
            <div className="-space-y-px shadow-2xl p-8 rounded-xl">
              <LoginForm onSuccess={onSuccess} />
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
