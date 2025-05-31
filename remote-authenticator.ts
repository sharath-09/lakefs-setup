import dotenv from "dotenv";
import express, { Express, Request, Response } from "express";
import { StatusCodes } from "http-status-codes";

type AuthRequestBody = {
  username: string;
  password: string;
};

type AuthResponseBody = {
  external_user_identifier: string;
};

const DEFAULT_PORT = 2080;

dotenv.config();

const port = process.env.PORT || DEFAULT_PORT;
const app: Express = express();

app.post(
  "/auth",
  (req: Request<AuthResponseBody, {}, AuthRequestBody>, res: Response) => {
    const { username, password } = req.body;
    if (!username?.length || !password?.length) {
      return res.status(StatusCodes.BAD_REQUEST).json({
        external_user_identifier: "",
      });
    }

    // 👇🏻 This is where you would implement your own authentication logic
    if (
      username === "testy.mctestface@example.com" &&
      password === "Password1"
    ) {
      return res.status(StatusCodes.OK).json({
        external_user_identifier: "TestyMcTestface",
      });
    } else {
      return res.status(StatusCodes.UNAUTHORIZED).json({
        external_user_identifier: "",
      });
    }
  }
);

app.listen(port, () => {
  console.log(`Remote Authenticator listening on port ${port}`);
});
