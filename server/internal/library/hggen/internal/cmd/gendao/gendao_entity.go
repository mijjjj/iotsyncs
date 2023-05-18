package gendao

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"hotgo/internal/library/hggen/internal/consts"
	"hotgo/internal/library/hggen/internal/utility/mlog"
	"hotgo/internal/library/hggen/internal/utility/utils"
)

func generateEntity(ctx context.Context, in CGenDaoInternalInput) {
	var dirPathEntity = gfile.Join(in.Path, in.EntityPath)
	if in.Clear {
		doClear(ctx, dirPathEntity, false)
	}
	// Model content.
	for i, tableName := range in.TableNames {
		fieldMap, err := in.DB.TableFields(ctx, tableName)
		if err != nil {
			mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", tableName, err)
		}
		var (
			newTableName   = in.NewTableNames[i]
			entityFilePath = gfile.Join(dirPathEntity, gstr.CaseSnake(newTableName)+".go")
			entityContent  = generateEntityContent(
				in,
				newTableName,
				gstr.CaseCamel(newTableName),
				generateStructDefinition(ctx, generateStructDefinitionInput{
					CGenDaoInternalInput: in,
					TableName:            tableName,
					StructName:           gstr.CaseCamel(newTableName),
					FieldMap:             fieldMap,
					IsDo:                 false,
				}),
			)
		)
		err = gfile.PutContents(entityFilePath, strings.TrimSpace(entityContent))
		if err != nil {
			mlog.Fatalf("writing content to '%s' failed: %v", entityFilePath, err)
		} else {
			utils.GoFmt(entityFilePath)
			mlog.Print("generated:", entityFilePath)
		}
	}
}

func generateEntityContent(in CGenDaoInternalInput, tableName, tableNameCamelCase, structDefine string) string {
	entityContent := gstr.ReplaceByMap(
		getTemplateFromPathOrDefault(in.TplDaoEntityPath, consts.TemplateGenDaoEntityContent),
		g.MapStrStr{
			tplVarTableName:          tableName,
			tplVarPackageImports:     getImportPartContent(structDefine, false),
			tplVarTableNameCamelCase: tableNameCamelCase,
			tplVarStructDefine:       structDefine,
		})
	entityContent = replaceDefaultVar(in, entityContent)
	return entityContent
}
