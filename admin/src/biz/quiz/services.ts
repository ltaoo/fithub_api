import { FetchParams } from "@/domains/list/typing";
import { TmpRequestResp } from "@/domains/request/utils";
import { Result } from "@/domains/result";
import { request } from "@/biz/requests";
import { ListResponse } from "@/biz/requests/types";

import { QuizTypes } from "./constants";

export function createQuiz(body: {
  content: string;
  overview: string;
  type: QuizTypes;
  difficulty: number;
  tags: string[];
  choices: { value: number; text: string }[];
  answer: { type: QuizTypes; value: number[] };
  analysis: string;
}) {
  return request.post("/api/quiz/create", {
    content: body.content,
    overview: body.overview,
    type: body.type,
    difficulty: body.difficulty,
    tags: body.tags.join(","),
    choices: JSON.stringify(body.choices),
    answer: JSON.stringify(body.answer),
    analysis: body.analysis,
  });
}

export function fetchQuizList(body: FetchParams) {
  return request.post<
    ListResponse<{
      id: number;
      content: string;
      choices: string;
    }>
  >("/api/quiz/list", {
    page_size: body.pageSize,
    page: body.page,
  });
}

export function createPaper(body: {
  name: string;
  overview: string;
  tags: string[];
  pass_score: number;
  duration: number;
  quiz_list: {
    id: number;
    score: number;
    sort_idx: number;
  }[];
}) {
  return request.post("/api/paper/create", {
    name: body.name,
    overview: body.overview,
    tags: body.tags.join(","),
    pass_score: body.pass_score,
    duration: body.duration,
    quiz_list: body.quiz_list,
  });
}

export function updatePaper(body: {
  id: number;
  name: string;
  overview: string;
  tags: string[];
  pass_score: number;
  duration: number;
  quiz_list: {
    relation_id?: number;
    id: number;
    score: number;
    sort_idx: number;
  }[];
}) {
  return request.post("/api/paper/update", {
    id: body.id,
    name: body.name,
    overview: body.overview,
    tags: body.tags.join(","),
    pass_score: body.pass_score,
    duration: body.duration,
    quiz_list: body.quiz_list,
  });
}

export function fetchPaperProfile(body: { id: number }) {
  return request.post<{
    paper: {
      id: number;
      name: string;
      overview: string;
      tags: string;
      pass_score: number;
      duration: number;
    };
    quizzes: {
      id?: number;
      quiz_id: number;
      quiz: {
        content: string;
      };
      score: number;
      sort_idx: number;
    }[];
  }>("/api/paper/profile", {
    id: body.id,
  });
}
export function fetchPaperProfileProcess(r: TmpRequestResp<typeof fetchPaperProfile>) {
  if (r.error) {
    return Result.Err(r.error.message);
  }
  const resp = r.data;
  return Result.Ok({
    id: resp.paper.id,
    name: resp.paper.name,
    overview: resp.paper.overview,
    tags: resp.paper.tags.split(","),
    pass_score: resp.paper.pass_score,
    duration: resp.paper.duration,
    quizzes: resp.quizzes.map((quiz) => {
      return {
        relation_id: quiz.id,
        id: quiz.quiz_id,
        content: quiz.quiz.content,
        score: quiz.score,
        sort_idx: quiz.sort_idx,
      };
    }),
  });
}

export function fetchPaperList(body: FetchParams) {
  return request.post<
    ListResponse<{
      id: number;
      name: string;
    }>
  >("/api/paper/list", {
    page_size: body.pageSize,
    page: body.page,
  });
}
