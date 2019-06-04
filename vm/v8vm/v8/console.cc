#include "console.h"
#include <iostream>

static consoleLogFunc CConsole = nullptr;
static sysLogFunc CSysLog = nullptr;

void InitGoConsole(consoleLogFunc console, sysLogFunc syslog) {
    CConsole = console;
    CSysLog = syslog;
}

void SysLog(SandboxPtr ptr, std::string level, std::string msg) {
    CSysLog(
        ptr,
        {const_cast<char *>(level.c_str()), static_cast<int>(level.length())},
        {const_cast<char *>(msg.c_str()), static_cast<int>(msg.length())}
    );
}

void NewConsoleLog(const FunctionCallbackInfo<Value> &args) {
    Isolate *isolate = args.GetIsolate();
    Local<Context> context = isolate->GetCurrentContext();
    Local<Object> global = context->Global();
    Local<Value> val = global->GetInternalField(0);
    if (!val->IsExternal()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "consoleLog val error")
        );
        isolate->ThrowException(err);
        return;
    }
    SandboxPtr sbxPtr = static_cast<SandboxPtr>(Local<External>::Cast(val)->Value());

    if (args.Length() != 2) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "consoleLog invalid argument length")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> levelVal = args[0];
    if (!levelVal->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "consoleLog log level must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    Local<Value> logVal = args[1];
    if (!logVal->IsString()) {
        Local<Value> err = Exception::Error(
            String::NewFromUtf8(isolate, "consoleLog log message must be string")
        );
        isolate->ThrowException(err);
        return;
    }

    NewCStrChecked(levelStr, levelVal, isolate);
    NewCStrChecked(logStr, logVal, isolate);

    CConsole(sbxPtr, levelStr, logStr);
}

void InitConsole(Isolate *isolate, Local<ObjectTemplate> globalTpl) {
    globalTpl->Set(
        String::NewFromUtf8(isolate, "_cLog", NewStringType::kNormal)
                    .ToLocalChecked(),
        FunctionTemplate::New(isolate, NewConsoleLog));
}
