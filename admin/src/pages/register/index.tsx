/**
 * @file 用户注册
 */
import { Button, Input, Card, CardHeader, CardTitle, CardContent, Label } from "@/components/ui";
import { InputCore, ButtonCore } from "@/domains/ui";
import { ViewComponent } from "@/store/types";

export const RegisterPage: ViewComponent = (props) => {
  const { app } = props;
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
    onEnter() { registerBtn.click(); },
    placeholder: "请输入密码",
    onChange(v) {
      app.$user.inputPassword(v);
    },
  });
  const registerBtn = new ButtonCore({
    async onClick() {
      registerBtn.setLoading(true);
      const r = await app.$user.register();
      registerBtn.setLoading(false);
      if (r.error) {
        app.tip({
          text: [r.error.message],
        });
        return;
      }
      app.tip({
        text: ["注册成功"],
      });
    },
  });

  return (
    <div class="flex min-h-screen items-center justify-center bg-muted/30 p-6">
      <Card class="w-full max-w-sm">
        <CardHeader><CardTitle class="text-xl">Fithub 教练注册</CardTitle><p class="text-sm text-muted-foreground">创建教练账户</p></CardHeader>
        <CardContent>
          <form class="space-y-4" onSubmit={(event) => { event.preventDefault(); registerBtn.click(); }}>
            <div class="space-y-2"><Label for="register-email">邮箱</Label><Input id="register-email" store={emailInput} /></div>
            <div class="space-y-2"><Label for="register-password">密码</Label><Input id="register-password" store={passwordInput} /></div>
            <Button store={registerBtn} class="w-full">注册</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};
