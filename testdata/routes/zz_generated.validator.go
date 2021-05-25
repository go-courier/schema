/*
Package routes GENERATED BY gengo:validator 
DON'T EDIT THIS FILE
*/
package routes

import (
	mime_multipart "mime/multipart"

	github_com_go_courier_schema_pkg_validator "github.com/go-courier/schema/pkg/validator"
	github_com_go_courier_schema_testdata_a "github.com/go-courier/schema/testdata/a"
)

func validateCookie(v *Cookie) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *Cookie) Validate() error {
	return validateCookie(v)
}
func validateCreate(v *Create) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateCreateFieldData(&v.Data), github_com_go_courier_schema_pkg_validator.Location("body"), "")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateCreateFieldData(v *Data) error {
	return validateData(v)
}
func validateData(v *Data) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateDataFieldData(&v.Data), "data")
	errSet.AddErr(validateDataFieldSubData(&v.SubData), "subData")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateDataFieldData(v **Data) error {
	vv := *v
	if vv == nil {
		return nil
	}
	return validateDataFieldDataPtr(vv)

}
func validateDataFieldDataPtr(v *Data) error {
	return validateData(v)
}
func validateDataFieldSubData(v **SubData) error {
	vv := *v
	if vv == nil {
		return nil
	}
	return validateDataFieldSubDataPtr(vv)

}
func validateDataFieldSubDataPtr(v *SubData) error {
	return validateSubData(v)
}
func validateSubData(v *SubData) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateCreateFieldDataFieldData(v **Data) error {
	vv := *v
	if vv == nil {
		return nil
	}
	return validateCreateFieldDataFieldDataPtr(vv)

}
func validateCreateFieldDataFieldDataPtr(v *Data) error {
	return validateData(v)
}
func validateCreateFieldDataFieldSubData(v **SubData) error {
	vv := *v
	if vv == nil {
		return nil
	}
	return validateCreateFieldDataFieldSubDataPtr(vv)

}
func validateCreateFieldDataFieldSubDataPtr(v *SubData) error {
	return validateSubData(v)
}
func (v *Create) Validate() error {
	return validateCreate(v)
}
func validateDataProvider(v *DataProvider) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateDataProviderFieldID(&v.ID), github_com_go_courier_schema_pkg_validator.Location("path"), "id")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateDataProviderFieldID(v *string) error {
	vv := *v
	if vv == "" {
		return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
	}
	return (&(github_com_go_courier_schema_pkg_validator.StringValidator{
		MinLength: 6,
	})).Validate(string(vv))

}
func (v *DataProvider) Validate() error {
	return validateDataProvider(v)
}
func validateDownloadFile(v *DownloadFile) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *DownloadFile) Validate() error {
	return validateDownloadFile(v)
}
func validateFormMultipartWithFile(v *FormMultipartWithFile) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateFormMultipartWithFileFieldFormData(&v.FormData), github_com_go_courier_schema_pkg_validator.Location("body"), "")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateFormMultipartWithFileFieldFormData(v *struct {
	String string                     `name:"string,omitempty"`
	Slice  []string                   `name:"slice,omitempty"`
	Data   Data                       `name:"data,omitempty"`
	File   *mime_multipart.FileHeader `name:"file" validate:"-"`
}) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateFormMultipartWithFileFieldFormDataFieldSlice(&v.Slice), "slice")
	errSet.AddErr(validateFormMultipartWithFileFieldFormDataFieldData(&v.Data), "data")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateFormMultipartWithFileFieldFormDataFieldSlice(v *[]string) error {
	list := *v
	n := len(list)
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for i := 0; i < n; i++ {
		errSet.AddErr(validateFormMultipartWithFileFieldFormDataFieldSliceElem(&list[i]), i)
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateFormMultipartWithFileFieldFormDataFieldSliceElem(v *string) error {
	return nil
}
func validateFormMultipartWithFileFieldFormDataFieldData(v *Data) error {
	return validateData(v)
}
func (v *FormMultipartWithFile) Validate() error {
	return validateFormMultipartWithFile(v)
}
func validateFormMultipartWithFiles(v *FormMultipartWithFiles) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateFormMultipartWithFilesFieldFormData(&v.FormData), github_com_go_courier_schema_pkg_validator.Location("body"), "")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateFormMultipartWithFilesFieldFormData(v *struct {
	Files []*mime_multipart.FileHeader `name:"files" validate:"-"`
}) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *FormMultipartWithFiles) Validate() error {
	return validateFormMultipartWithFiles(v)
}
func validateFormURLEncoded(v *FormURLEncoded) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateFormURLEncodedFieldFormData(&v.FormData), github_com_go_courier_schema_pkg_validator.Location("body"), "")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateFormURLEncodedFieldFormData(v *struct {
	String string   `name:"string"`
	Slice  []string `name:"slice"`
	Data   Data     `name:"data"`
}) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateFormURLEncodedFieldFormDataFieldSlice(&v.Slice), "slice")
	errSet.AddErr(validateFormURLEncodedFieldFormDataFieldData(&v.Data), "data")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateFormURLEncodedFieldFormDataFieldSlice(v *[]string) error {
	list := *v
	n := len(list)
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for i := 0; i < n; i++ {
		errSet.AddErr(validateFormURLEncodedFieldFormDataFieldSliceElem(&list[i]), i)
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateFormURLEncodedFieldFormDataFieldSliceElem(v *string) error {
	return nil
}
func validateFormURLEncodedFieldFormDataFieldData(v *Data) error {
	return validateData(v)
}
func (v *FormURLEncoded) Validate() error {
	return validateFormURLEncoded(v)
}
func validateGetByID(v *GetByID) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateGetByIDFieldProtocol(&v.Protocol), github_com_go_courier_schema_pkg_validator.Location("query"), "protocol")
	errSet.AddErr(validateGetByIDFieldBytes(&v.Bytes), github_com_go_courier_schema_pkg_validator.Location("query"), "bytes")
	errSet.AddErr(validateGetByIDFieldLabel(&v.Label), github_com_go_courier_schema_pkg_validator.Location("query"), "label")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateGetByIDFieldProtocol(v *[]github_com_go_courier_schema_testdata_a.Protocol) error {
	list := *v
	n := len(list)
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for i := 0; i < n; i++ {
		errSet.AddErr(validateGetByIDFieldProtocolElem(&list[i]), i)
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateGetByIDFieldProtocolElem(v *github_com_go_courier_schema_testdata_a.Protocol) error {
	return nil
}
func validateGetByIDFieldBytes(v *[]uint8) error {
	list := *v
	n := len(list)
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for i := 0; i < n; i++ {
		errSet.AddErr(validateGetByIDFieldBytesElem(&list[i]), i)
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateGetByIDFieldBytesElem(v *uint8) error {
	return nil
}
func validateGetByIDFieldLabel(v *[]string) error {
	list := *v
	n := len(list)
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	for i := 0; i < n; i++ {
		errSet.AddErr(validateGetByIDFieldLabelElem(&list[i]), i)
	}
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateGetByIDFieldLabelElem(v *string) error {
	return nil
}
func (v *GetByID) Validate() error {
	return validateGetByID(v)
}
func validateGetByUser(v *GetByUser) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *GetByUser) Validate() error {
	return validateGetByUser(v)
}
func validateHealthCheck(v *HealthCheck) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *HealthCheck) Validate() error {
	return validateHealthCheck(v)
}
func validateMustValidAccount(v *MustValidAccount) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *MustValidAccount) Validate() error {
	return validateMustValidAccount(v)
}
func validateProxy(v *Proxy) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *Proxy) Validate() error {
	return validateProxy(v)
}
func validateProxyV2(v *ProxyV2) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *ProxyV2) Validate() error {
	return validateProxyV2(v)
}
func validateRedirect(v *Redirect) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *Redirect) Validate() error {
	return validateRedirect(v)
}
func validateRedirectWhenError(v *RedirectWhenError) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *RedirectWhenError) Validate() error {
	return validateRedirectWhenError(v)
}
func validateRemoveByID(v *RemoveByID) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *RemoveByID) Validate() error {
	return validateRemoveByID(v)
}
func validateShowImage(v *ShowImage) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func (v *ShowImage) Validate() error {
	return validateShowImage(v)
}
func validateUpdateByID(v *UpdateByID) error {
	errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
	errSet.AddErr(validateUpdateByIDFieldData(&v.Data), github_com_go_courier_schema_pkg_validator.Location("body"), "")
	if errSet.Len() > 0 {
		return errSet
	}
	return nil

}
func validateUpdateByIDFieldData(v *Data) error {
	return validateData(v)
}
func (v *UpdateByID) Validate() error {
	return validateUpdateByID(v)
}
