import { PageKeysType, build } from "@/domains/route_view/utils";

/**
 * @file 路由配置
 */
const configure = {
  root: {
    title: "ROOT",
    pathname: "/",
    children: {
      home_layout: {
        title: "首页布局",
        pathname: "/home",
        children: {
          index: {
            title: "首页",
            pathname: "/home/index",
            options: {
              require: ["login"],
            },
          },
          muscle: {
            title: "肌肉列表",
            pathname: "/home/muscle",
            options: {
              require: ["login"],
            },
          },
          equipment: {
            title: "器械列表",
            pathname: "/home/equipment",
            options: {
              require: ["login"],
            },
          },
          action_list: {
            title: "动作列表",
            pathname: "/home/actions",
            options: {
              require: ["login"],
            },
          },
          action_create: {
            title: "动作创建",
            pathname: "/home/action_create",
            options: {
              require: ["login"],
            },
          },
          action_update: {
            title: "动作编辑",
            pathname: "/home/action_update",
            options: {
              require: ["login"],
            },
          },
          workout_plan_list: {
            title: "训练计划列表",
            pathname: "/home/workout_plans",
            options: {
              require: ["login"],
            },
          },
          workout_plan_create: {
            title: "训练计划创建",
            pathname: "/home/workout_plan_create",
            options: {
              require: ["login"],
            },
          },
          workout_plan_update: {
            title: "训练计划编辑",
            pathname: "/home/workout_plan_update",
            options: {
              require: ["login"],
            },
          },
          workout_day: {
            title: "训练日",
            pathname: "/home/workout_day",
            options: {
              require: ["login"],
            },
          },
          student_list: {
            title: "学员列表",
            pathname: "/home/students",
            options: {
              require: ["login"],
            },
          },
          student_create: {
            title: "学员创建",
            pathname: "/home/student_create",
            options: {
              require: ["login"],
            },
          },
          student_update: {
            title: "学员编辑",
            pathname: "/home/student_update",
            options: {
              require: ["login"],
            },
          },
          student_profile: {
            title: "学员详情",
            pathname: "/home/student_profile",
            options: {
              require: ["login"],
            },
          },
          plan_set_list: {
            title: "训练计划集合列表",
            pathname: "/home/plan_set/list",
            options: {
              require: ["login"],
            },
          },
          plan_set_create: {
            title: "创建训练计划集合",
            pathname: "/home/plan_set/create",
            options: {
              require: ["login"],
            },
          },
          plan_set_update: {
            title: "更新训练计划集合",
            pathname: "/home/plan_set/update",
            options: {
              require: ["login"],
            },
          },
          subscription_plan_create: {
            title: "创建订阅计划",
            pathname: "/home/subscription_plan/create",
            options: {
              require: ["login"],
            },
          },
          subscription_plan_list: {
            title: "订阅计划列表",
            pathname: "/home/subscription_plan/list",
            options: {
              require: ["login"],
            },
          },
          quiz_list: {
            title: "题目列表",
            pathname: "/home/quiz/list",
            options: {
              require: ["login"],
            },
          },
          quiz_create: {
            title: "创建题目",
            pathname: "/home/quiz/create",
            options: {
              require: ["login"],
            },
          },
          paper_list: {
            title: "试卷列表",
            pathname: "/home/paper/list",
            options: {
              require: ["login"],
            },
          },
          paper_create: {
            title: "创建试卷",
            pathname: "/home/paper/create",
            options: {
              require: ["login"],
            },
          },
          paper_update: {
            title: "编辑试卷",
            pathname: "/home/paper/update",
            options: {
              require: ["login"],
            },
          },
          gift_card_list: {
            title: "礼品卡",
            pathname: "/home/gift_card/list",
            options: {
              require: ["login"],
            },
          },
          gift_card_create: {
            title: "创建礼品卡",
            pathname: "/home/gift_card/create",
            options: {
              require: ["login"],
            },
          },
          gift_card_reward_create: {
            title: "创建奖励",
            pathname: "/home/gift_card_reward/create",
            options: {
              require: ["login"],
            },
          },
          coach_list: {
            title: "用户管理",
            pathname: "/home/coach/list",
            options: {
              require: ["login"],
            },
          },
          coach_create: {
            title: "创建用户",
            pathname: "/home/coach/create",
            options: {
              require: ["login"],
            },
          },
          coach_content_create: {
            title: "创建奖励",
            pathname: "/home/coach/profile",
            options: {
              require: ["login"],
            },
          },
          // content_of_workout_action: {
          //   title: "关联动作",
          //   pathname: "/home/workout_action/content/create",
          //   options: {
          //     require: ["login"],
          //   },
          // },
          content_of_workout_plan_create: {
            title: "关联计划",
            pathname: "/home/workout_plan/content/create",
            options: {
              require: ["login"],
            },
          },
        },
      },
      login: {
        title: "教练登录",
        pathname: "/login",
      },
      register: {
        title: "教练注册",
        pathname: "/register",
      },
      notfound: {
        title: "404",
        pathname: "/notfound",
      },
    },
  },
};
export type PageKeys = PageKeysType<typeof configure>;
const result = build<PageKeys>(configure);
export const routes = result.routes;
export const routesWithPathname = result.routesWithPathname;
