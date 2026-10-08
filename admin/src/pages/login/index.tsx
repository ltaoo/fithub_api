/**
 * @file 用户登录
 */
import { Button, Input, Card, CardHeader, CardTitle, CardContent, Label } from "@/components/ui";
import { ButtonCore, InputCore } from "@/domains/ui";
import { ViewComponent } from "@/store/types";

export const LoginPage: ViewComponent = (props) => {
  const { app, view } = props;
  const emailInput = new InputCore({
    defaultValue: "",
    placeholder: "请输入邮箱",
    onChange(v) {
      app.$user.inputEmail(v);
    },
  });
  const passwordInput = new InputCore({
    defaultValue: "",
    type: "password",
    onEnter() { loginBtn.click(); },
    placeholder: "请输入密码",
    onChange(v) {
      app.$user.inputPassword(v);
    },
  });
  const loginBtn = new ButtonCore({
    async onClick() {
      loginBtn.setLoading(true);
      await app.$user.login();
      loginBtn.setLoading(false);
    },
  });

  // view.onShow(() => {
  //   if (!app.user.needRegister) {
  //     return;
  //   }
  // });

  return (
    <div class="flex min-h-screen items-center justify-center bg-muted/30 p-6">
      <Card class="w-full max-w-sm">
        <CardHeader><CardTitle class="text-xl">Fithub 登录</CardTitle><p class="text-sm text-muted-foreground">登录健身管理后台</p></CardHeader>
        <CardContent>
          <form class="space-y-4" onSubmit={(event) => { event.preventDefault(); loginBtn.click(); }}>
            <div class="space-y-2"><Label for="login-email">邮箱</Label><Input id="login-email" store={emailInput} /></div>
            <div class="space-y-2"><Label for="login-password">密码</Label><Input id="login-password" store={passwordInput} /></div>
            <Button store={loginBtn} class="w-full">登录</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};
