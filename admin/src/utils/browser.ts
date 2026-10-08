import { Result } from "@/domains/result";

export function loadImage(data: any): Promise<Result<HTMLImageElement>> {
  return new Promise((resolve) => {
    const img = new Image();
    img.onload = () => {
      resolve(Result.Ok(img));
    };
    img.onerror = (msg) => {
      resolve(Result.Err(msg as string));
    };
    img.src = data;
  });
}

export function blobToArrayBuffer(blob: Blob): Promise<ArrayBuffer> {
  return new Promise((resolve) => {
    const reader = new FileReader();
    reader.onload = () => {
      const buffer = reader.result;
      resolve(buffer as ArrayBuffer);
    };
    reader.readAsArrayBuffer(blob);
  });
}

export function readFileAsURL(file: File): Promise<Result<string>> {
  return new Promise((resolve) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      if (e.target === null) {
        return resolve(Result.Err("read failed"));
      }
      resolve(Result.Ok(reader.result as string));
    };
    reader.readAsDataURL(file);
  });
}

export function readFileAsText(file: File): Promise<Result<string>> {
  const reader = new FileReader();
  return new Promise((resolve) => {
    reader.onload = (e) => {
      try {
        // const jsonContent = JSON.parse(e.target?.result as string);
        // 这里可以处理解析后的 JSON 内容
        // console.log("JSON content:", jsonContent);
        if (!e.target) {
          resolve(Result.Err("读取失败"));
          return;
        }
        resolve(Result.Ok(e.target.result as string));
      } catch (error) {
        const e = error as Error;
        resolve(Result.Err(e.message));
      }
    };
    reader.readAsText(file);
  });
}

export function readFileAsArrayBuffer(file: File): Promise<Result<ArrayBuffer>> {
  return new Promise((resolve) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      if (e.target === null) {
        return resolve(Result.Err("read failed"));
      }
      const buffer = e.target.result;
      return resolve(Result.Ok(buffer as ArrayBuffer));
    };
    reader.readAsArrayBuffer(file);
  });
}
