import { JSXElement } from "@/timeless";

import { ViewComponent } from "@/store/types";

import { PageKeys } from "./routes";

import { LoginPage } from "@/pages/login";
import { RegisterPage } from "@/pages/register";
import { NotFoundPage } from "@/pages/notfound";
import { HomeLayout } from "@/pages/home/layout";
import { HomeIndexPage } from "@/pages/home";
import { WorkoutActionListView } from "@/pages/workout_action/list";
import { WorkoutActionCreateView } from "@/pages/workout_action/create";
import { WorkoutActionUpdateView } from "@/pages/workout_action/update";
import { WorkoutPlanListView } from "@/pages/workout_plan/list";
import { WorkoutPlanCreateView } from "@/pages/workout_plan/create";
import { WorkoutPlanUpdatePage } from "@/pages/workout_plan/update";
import { StudentListView } from "@/pages/student/list";
import { StudentCreateView } from "@/pages/student/create";
import { StudentUpdateView } from "@/pages/student/update";
import { StudentProfileView } from "@/pages/student/profile";
import { MuscleListView } from "@/pages/muscle/list";
import { EquipmentListView } from "@/pages/equipment/list";
import { WorkoutDayUpdateView } from "@/pages/workout_day/create";
import { WorkoutPlanSetListView } from "@/pages/workout_plan_set/list";
import { WorkoutPlanSetUpdateView } from "@/pages/workout_plan_set/update";
import { WorkoutPlanSetCreateView } from "@/pages/workout_plan_set/create";
import { SubscriptionPlanCreateView } from "@/pages/subscription_plan/create";
import { SubscriptionPlanListView } from "@/pages/subscription_plan/list";
import { PaperCreateView } from "@/pages/paper/create";
import { PaperListView } from "@/pages/paper/list";
import { QuizListView } from "@/pages/quiz/list";
import { QuizCreateView } from "@/pages/quiz/create";
import { PaperUpdateView } from "@/pages/paper/update";
import { GiftCardRewardCreateView } from "@/pages/gift_card_reward/create";
import { GiftCardListView } from "@/pages/gift_card/list";
import { GiftCardCreateView } from "@/pages/gift_card/create";
import { InfluencerCreateView } from "@/pages/coach/create";
import { InfluencerContentUploadView } from "@/pages/coach/upload_content";
import { CoachListView } from "@/pages/coach/list";
import { ContentOfWorkoutPlanCreateView } from "@/pages/workout_plan/link_content_create";

export const pages: Omit<Record<PageKeys, ViewComponent>, "root"> = {
  "root.login": LoginPage,
  "root.register": RegisterPage,
  "root.notfound": NotFoundPage,
  "root.home_layout": HomeLayout,
  "root.home_layout.index": HomeIndexPage,
  "root.home_layout.muscle": MuscleListView,
  "root.home_layout.equipment": EquipmentListView,
  "root.home_layout.action_list": WorkoutActionListView,
  "root.home_layout.action_create": WorkoutActionCreateView,
  "root.home_layout.action_update": WorkoutActionUpdateView,
  "root.home_layout.workout_plan_list": WorkoutPlanListView,
  "root.home_layout.workout_plan_create": WorkoutPlanCreateView,
  "root.home_layout.workout_plan_update": WorkoutPlanUpdatePage,
  "root.home_layout.workout_day": WorkoutDayUpdateView,
  "root.home_layout.student_list": StudentListView,
  "root.home_layout.student_create": StudentCreateView,
  "root.home_layout.student_update": StudentUpdateView,
  "root.home_layout.student_profile": StudentProfileView,
  "root.home_layout.plan_set_list": WorkoutPlanSetListView,
  "root.home_layout.plan_set_update": WorkoutPlanSetUpdateView,
  "root.home_layout.plan_set_create": WorkoutPlanSetCreateView,
  "root.home_layout.subscription_plan_list": SubscriptionPlanListView,
  "root.home_layout.subscription_plan_create": SubscriptionPlanCreateView,
  "root.home_layout.paper_list": PaperListView,
  "root.home_layout.paper_create": PaperCreateView,
  "root.home_layout.paper_update": PaperUpdateView,
  "root.home_layout.quiz_list": QuizListView,
  "root.home_layout.quiz_create": QuizCreateView,
  "root.home_layout.gift_card_list": GiftCardListView,
  "root.home_layout.gift_card_reward_create": GiftCardRewardCreateView,
  "root.home_layout.gift_card_create": GiftCardCreateView,
  "root.home_layout.coach_list": CoachListView,
  "root.home_layout.coach_create": InfluencerCreateView,
  "root.home_layout.coach_content_create": InfluencerContentUploadView,
  "root.home_layout.content_of_workout_plan_create": ContentOfWorkoutPlanCreateView,
};
