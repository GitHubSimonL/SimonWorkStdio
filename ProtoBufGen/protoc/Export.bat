@echo on
if exist export (
	rd /s /q export > NUL
)
%XML2PROTOC_PATH%\xml2proto.exe go cproto.xml
%XML2PROTOC_PATH%\xml2proto.exe go sproto.xml
mkdir export
%PROTOC_PATH%\protoc.exe --go_out=export protocol\common\common.proto
rem %PROTOC_PATH%\protoc.exe --go_out=export Arkenstone.proto
%PROTOC_PATH%\protoc.exe --go_out=plugins=grpc:.\export\protocol CProto.proto
%PROTOC_PATH%\protoc.exe --go_out=plugins=grpc:.\export\protocol SProto.proto

..\..\tools\protoc\protoPostHandle.exe -src export\protocol\common\common.pb.go -descName commonDescriptor -out export\protocol\common\common.pb.go
..\..\tools\protoc\protoPostHandle.exe -src export\protocol\CProto.pb.go -descName cprotoDescriptor -out export\protocol\CProto.pb.go
..\..\tools\protoc\protoPostHandle.exe -src export\protocol\SProto.pb.go -descName sprotoDescriptor -out export\protocol\SProto.pb.go

%PostCommand_PATH%\postCommand.exe export\protocol\CProto.pb.go
%PostConvertProtobufStructs_PATH%\postConvertProtobufStructs.exe
..\..\tools\protoc\protoc.exe --csharp_out=export protocol\common\common.proto
..\..\tools\protoc\protoc.exe --csharp_out=export CProto.proto
rem del .\*.go
del .\sproto.proto
del .\cproto.proto
%XML2PROTOC_PATH%\xml2proto.exe lua cproto.xml

%PROTOC_PATH%\protoc.exe protocol\common\common.proto cproto.proto  -o .\cproto.png
python .\makePBHeader.py
del .\cproto.png
copy cproto.h ..\..\client\trunk\src\cproto.h
del .\cproto.h
copy cproto.proto ..\..\client\trunk\data\proto\cproto.proto

..\..\client\trunk\tools\luatools\lua.exe ..\..\tools\makeLuaAPI\clientMakeLuaAPI.lua

copy export\protocol\common\common.pb.go %PROJ_DIR%\server\lichen\src\protocol\common
copy export\protocol\CProto.pb.go %PROJ_DIR%\server\lichen\src\protocol\
copy export\protocol\SProto.pb.go %PROJ_DIR%\server\lichen\src\protocol\
rem copy export\protocol\common\common.pb.go W:\server\lichen\src\protocol\common
rem copy export\protocol\CProto.pb.go W:\server\lichen\src\protocol\
rem copy export\protocol\SProto.pb.go W:\server\lichen\src\protocol\
