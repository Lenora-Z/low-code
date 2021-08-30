// Package server
//Created by Goland
//@User: lenora
//@Date: 2021/2/5
//@Time: 2:40 下午
package server

func (ds *defaultServer) routers() {
	platformAuth := ds.tritiumAuth()
	appAuth := ds.appAuth()
	application := ds.engine.Group("/app/")
	application.Use(platformAuth)
	{
		application.GET("list", ds.ApplicationList)
		application.POST("new", ds.NewApplication)
		application.POST("edit", ds.EditApplication)
		application.Use(appAuth)
		application.GET("detail", ds.ApplicationDetail)
	}

	form := ds.engine.Group("/form/")
	form.Use(platformAuth, appAuth)
	{
		form.GET("list", ds.GetFormList)
		//form.GET("list/main", ds.GetMainFormList)
		form.POST("new", ds.NewForm)
		form.POST("edit", ds.EditForm)
		//form.POST("copy", ds.CopyForm)
		form.GET("detail", ds.FormDetail)
		form.GET("fields", ds.FieldList)
		//form.GET("service/fields", ds.FieldsForService)
		form.GET("relation", ds.RelatedForm)
		form.GET("list/collect", ds.CollectFormList)
		form.POST("delete", ds.DeleteForm)
		form.GET("flow/button", ds.FlowButton)
		form.GET("mapping", ds.FlowMapping)
		form.GET("table/columns", ds.TablaFieldColumns)
	}

	nav := ds.engine.Group("/nav/")
	nav.Use(platformAuth, appAuth)
	{
		//nav.GET("list", ds.NavigationList)
		nav.POST("edit", ds.EditNavigation)
		nav.GET("detail", ds.RouteList)
	}

	route := ds.engine.Group("/route/")
	route.Use(platformAuth, appAuth)
	{
		route.GET("list", ds.RouteList)
	}

	flow := ds.engine.Group("/flow/")
	flow.Use(platformAuth, appAuth)
	{
		flow.GET("list", ds.FlowList)
		flow.POST("new", ds.NewFlow)
		flow.POST("edit", ds.EditFlow)
		flow.GET("detail", ds.FlowDetail)
		flow.GET("record", ds.FlowLogs)
		flow.GET("record/log", ds.FlowLogsDetail)
		flow.POST("delete", ds.DeleteFlow)
	}

	maps := ds.engine.Group("/map/")
	maps.Use(platformAuth, appAuth)
	{
		maps.POST("bind", ds.FlowBound)
		maps.GET("list", ds.OnlineMaps)
	}

	service := ds.engine.Group("/service/")
	service.Use(platformAuth, appAuth)
	{
		service.GET("list/all", ds.ServiceList)
		//service.GET("list", ds.AppService)
		service.POST("param/bind", ds.NewParamRelies)
		service.GET("params", ds.ServiceParams)
		service.GET("package", ds.PackageService)
	}

	group := ds.engine.Group("/organization/")
	group.Use(platformAuth, appAuth)
	{
		group.GET("list", ds.OrganizationList)
		group.POST("new", ds.newOrganization)
		group.POST("edit", ds.EditOrganization)
		group.POST("delete", ds.DeleteOrganization)
		group.POST("config/role", ds.ConfigOrgRole)
		group.GET("detail", ds.OrganizationDetail)
	}

	role := ds.engine.Group("/role/")
	role.Use(platformAuth, appAuth)
	{
		role.GET("list", ds.RoleList)
		role.POST("new", ds.NewRole)
		role.POST("edit", ds.EditRole)
		role.POST("delete", ds.DeleteRole)
		role.GET("detail", ds.RoleDetail)
	}

	user := ds.engine.Group("/user/")
	user.Use(platformAuth, appAuth)
	{
		user.POST("list", ds.UserList)
		user.GET("all", ds.AllUserList)
		user.POST("new", ds.NewUser)
		user.POST("edit", ds.EditUser)
		user.POST("password/reset", ds.PasswordReset)
		user.POST("delete", ds.DeleteUser)
		user.GET("detail", ds.UserDetail)
	}

	version := ds.engine.Group("/version/")
	version.Use(platformAuth, appAuth)
	{
		version.GET("list", ds.VersionHistory)
		version.POST("publish", ds.PublishApplication)
		version.GET("online", ds.OnlineVersion)
		version.POST("offline", ds.OfflineApplication)
	}

	file := ds.engine.Group("/file/")
	{
		file.GET("show", ds.FilePreView)
		file.Use(platformAuth, appAuth)
		file.POST("upload", ds.FileUpload)
	}

	table := ds.engine.Group("/table/")
	table.Use(platformAuth, appAuth)
	{
		table.GET("page", ds.TableListByPage)
		table.GET("columns", ds.TableColumns)
		table.GET("relations", ds.TableRelations)
		table.GET("relations/attr", ds.TableRelationsAttr)
		table.GET("model", ds.TableModel)
		table.GET("list", ds.TableList)
		table.GET("relation/columns", ds.TableRelationColumns)
		table.GET("columns/type", ds.TypeColumns)
		table.GET("meta", ds.TableMeta)
	}

	userAuth := ds.CheckToken()
	api := ds.engine.Group("/api/")
	{
		//api.POST("sms/send", ds.SendSms)

		//转发请求访问第三方接口
		api.Any("interface/call", ds.CallInterface)

		api.POST("hash", ds.ApplicationByHash)
		userNone := api.Group("user/")
		userNone.POST("login", ds.UserLogin)
		api.Use(userAuth)

		api.POST("flow/run", ds.CallBpmn)
		{
			api.GET("nav", ds.AppNavigation)
			api.POST("file/upload", ds.ClientFileUpload)
		}
		userNone.Use(userAuth)
		{
			userNone.POST("password/reset", ds.UserChangePwd)
			userNone.GET("info", ds.UserInfo)
			userNone.GET("permission", ds.UserPermission)
			userNone.GET("list", ds.ClientUserList)
		}

		orgNone := api.Group("organization/")
		{
			orgNone.GET("list", ds.ClientOrganizationList)
		}

		applyNone := api.Group("apply/")
		{
			applyNone.GET("list", ds.ApplyList)
			applyNone.GET("param", ds.ApplyParams)
			applyNone.POST("deal", ds.DealApply)
			applyNone.GET("status", ds.ApplyStatus)
			applyNone.GET("list/notifier", ds.NotifierApplyList)
			applyNone.GET("notifier/param", ds.ApplyParams)
		}

		formNone := api.Group("form")
		formAuth := ds.checkUserApiPermission()
		{
			formNone.GET("", ds.GetFormDetail)
			formNone.Use(formAuth)
			//formNone.GET("/detail", ds.FormDataDetail)
			//formNone.POST("/submit", ds.AddFormData)
			//formNone.POST("/edit", ds.EditFormData)
			//formNone.POST("/data", ds.FormDataList)
			//formNone.POST("/log", ds.FormDataLog)
			//formNone.POST("/delete", ds.DeleteFormData)
			//formNone.GET("/download", ds.FormDataListExport)
			//formNone.GET("/related", ds.RelatedFormData)
			formNone.POST("/cascade", ds.GetCascadeData)
			formNone.POST("/submit", ds.SubmitFormData)

			//列表控件相关
			tableNone := formNone.Group("/table")
			{
				tableNone.POST("/list", ds.GetTableDataByFilter)
				tableNone.POST("/format", ds.FormatTableData)
				tableNone.POST("/edit", ds.UpdateTableData)
				tableNone.DELETE("/delete", ds.DeleteTableData)
				tableNone.POST("/detail", ds.TableDataDetail)
				tableNone.GET("/option/list", ds.GetTableOptionData)
				tableNone.POST("/download", ds.ExportTableData)
			}

			formNone.POST("/relate", ds.GetRelatedData)
			formNone.POST("/relate/detail", ds.GetRelatedDataDetail)
		}
	}

	hashAuth := ds.appHashAuth()
	public := ds.engine.Group("/public/")
	{
		public.POST("data/new", ds.CreateData)
		public.POST("data/edit", ds.UpdateData)
		public.POST("package/sign", ds.Signature)
		public.POST("package/email", ds.SendEmail)
		public.POST("package/filing", ds.ContractFiling)
		public.Use(hashAuth)
		public.POST("form/data", ds.publicFormDataList)
	}

	hrm := ds.engine.Group("/hrm/")
	{
		hrm.GET("statistic", ds.GetHrmStatistics)
	}
}
