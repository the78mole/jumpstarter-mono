// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from "vscode";
import { PythonExtension } from "@vscode/python-extension";

const { execa } = require("execa");

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export async function activate(context: vscode.ExtensionContext): Promise<any> {
  // Use the console to output diagnostic information (console.log) and errors (console.error)
  // This line of code will only be executed once when your extension is activated
  console.log('Congratulations, your extension "jumpstarter" is now active!');

  // The command has been defined in the package.json file
  // Now provide the implementation of the command with registerCommand
  // The commandId parameter must match the command field in package.json
  const disposable = vscode.commands.registerCommand(
    "jumpstarter.helloWorld",
    () => {
      // The code you place here will be executed every time your command is executed
      // Display a message box to the user
      vscode.window.showInformationMessage(
        "Hello World from vscode-jumpstarter!",
      );
    },
  );

  const python = await getPythonExtension();
  if (python) {
    await checkJumpstarterVersion(python);
    const pythonEnvironment =
      python.environments.onDidChangeActiveEnvironmentPath((event) => {
        console.log("Python environment changed", event);
        checkJumpstarterVersion(python);
      });
    context.subscriptions.push(pythonEnvironment);
  }

  context.subscriptions.push(disposable);
}

async function checkJumpstarterVersion(
  python: PythonExtension,
): Promise<boolean> {
  const path = python.environments.getActiveEnvironmentPath().path;
  console.log("Found Python extension! Active environment path:", path);
  const { stdout } =
    await execa`${path} -m jumpstarter_cli version --output json`;
  console.log("Command Output: ", stdout);
  return true;
}

/**
 * Return the python extension's API, if available.
 */
async function getPythonExtension(): Promise<PythonExtension | undefined> {
  try {
    return await PythonExtension.api();
  } catch (err) {
    console.error(`Unable to load python extension: ${err}`);
    return undefined;
  }
}

// This method is called when your extension is deactivated
export function deactivate() {}
