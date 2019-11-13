package permission

//ValidatePermissions handlers middleware
// func ValidatePermissions(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		//get role_id from JWT
// 		user := c.Get("user")
// 		token := user.(*jwt.Token)
// 		claims := token.Claims.(jwt.MapClaims)
// 		RoleID := claims["role_id"]

// 		//method and url from request
// 		Method := c.Request().Method
// 		URL := c.Request().URL.String()
// 		PermissionsModel := models.Permissions{}

// 		if !asira.App.DB.Where("lower(permissions) = 'all' AND role_id = ?", RoleID).Find(&PermissionsModel).RecordNotFound() {
// 			return next(c)
// 		}
// 		//check permissions
// 		perConfig := asira.App.Permission.GetStringMap(fmt.Sprintf("%s", Method))
// 		validate, key := ifExist(perConfig, URL)
// 		if validate {
// 			if !asira.App.DB.Where("lower(permissions) = ? AND role_id = ?", key, RoleID).Find(&PermissionsModel).RecordNotFound() {
// 				return next(c)
// 			} else {
// 				return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "you are not allowed"))
// 			}
// 		} else {
// 			return next(c)
// 		}
// 	}
// }

// func ifExist(list map[string]interface{}, URL string) (bool, string) {
// 	for key, value := range list {
// 		if strings.Contains(URL, value.(string)) {
// 			return true, key
// 		}
// 	}
// 	return false, ""
// }
