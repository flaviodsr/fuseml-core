// Code generated by goa v3.3.1, DO NOT EDIT.
//
// runnable gRPC server types
//
// Command:
// $ goa gen github.com/fuseml/fuseml-core/design

package server

import (
	runnablepb "github.com/fuseml/fuseml-core/gen/grpc/runnable/pb"
	runnable "github.com/fuseml/fuseml-core/gen/runnable"
	goa "goa.design/goa/v3/pkg"
)

// NewListPayload builds the payload of the "list" endpoint of the "runnable"
// service from the gRPC request type.
func NewListPayload(message *runnablepb.ListRequest) *runnable.ListPayload {
	v := &runnable.ListPayload{}
	if message.Id != "" {
		v.ID = &message.Id
	}
	return v
}

// NewListResponse builds the gRPC response type from the result of the "list"
// endpoint of the "runnable" service.
func NewListResponse(result []*runnable.Runnable) *runnablepb.ListResponse {
	message := &runnablepb.ListResponse{}
	message.Field = make([]*runnablepb.Runnable2, len(result))
	for i, val := range result {
		message.Field[i] = &runnablepb.Runnable2{
			Name: val.Name,
			Kind: val.Kind,
		}
		if val.ID != nil {
			message.Field[i].Id = *val.ID
		}
		if val.Created != nil {
			message.Field[i].Created = *val.Created
		}
		if val.Image != nil {
			message.Field[i].Image = svcRunnableRunnableImageToRunnablepbRunnableImage(val.Image)
		}
		if val.Inputs != nil {
			message.Field[i].Inputs = make([]*runnablepb.RunnableInput, len(val.Inputs))
			for j, val := range val.Inputs {
				message.Field[i].Inputs[j] = &runnablepb.RunnableInput{}
				if val.Name != nil {
					message.Field[i].Inputs[j].Name = *val.Name
				}
				if val.Kind != nil {
					message.Field[i].Inputs[j].Kind = *val.Kind
				}
				if val.Runnable != nil {
					message.Field[i].Inputs[j].Runnable = *val.Runnable
				}
			}
		}
		if val.Outputs != nil {
			message.Field[i].Outputs = make([]*runnablepb.RunnableOutput, len(val.Outputs))
			for j, val := range val.Outputs {
				message.Field[i].Outputs[j] = &runnablepb.RunnableOutput{}
				if val.Name != nil {
					message.Field[i].Outputs[j].Name = *val.Name
				}
				if val.Kind != nil {
					message.Field[i].Outputs[j].Kind = *val.Kind
				}
				if val.Runnable != nil {
					message.Field[i].Outputs[j].Runnable = svcRunnableRunnableRefToRunnablepbRunnableRef(val.Runnable)
				}
				if val.Metadata != nil {
					message.Field[i].Outputs[j].Metadata = svcRunnableInputParameterToRunnablepbInputParameter(val.Metadata)
				}
			}
		}
		if val.Labels != nil {
			message.Field[i].Labels = make([]string, len(val.Labels))
			for j, val := range val.Labels {
				message.Field[i].Labels[j] = val
			}
		}
	}
	return message
}

// NewRegisterPayload builds the payload of the "register" endpoint of the
// "runnable" service from the gRPC request type.
func NewRegisterPayload(message *runnablepb.RegisterRequest) *runnable.Runnable {
	v := &runnable.Runnable{
		Name: message.Name,
		Kind: message.Kind,
	}
	if message.Id != "" {
		v.ID = &message.Id
	}
	if message.Created != "" {
		v.Created = &message.Created
	}
	if message.Image != nil {
		v.Image = protobufRunnablepbRunnableImageToRunnableRunnableImage(message.Image)
	}
	if message.Inputs != nil {
		v.Inputs = make([]*runnable.RunnableInput, len(message.Inputs))
		for i, val := range message.Inputs {
			v.Inputs[i] = &runnable.RunnableInput{}
			if val.Name != "" {
				v.Inputs[i].Name = &val.Name
			}
			if val.Kind != "" {
				v.Inputs[i].Kind = &val.Kind
			}
			if val.Runnable != "" {
				v.Inputs[i].Runnable = &val.Runnable
			}
		}
	}
	if message.Outputs != nil {
		v.Outputs = make([]*runnable.RunnableOutput, len(message.Outputs))
		for i, val := range message.Outputs {
			v.Outputs[i] = &runnable.RunnableOutput{}
			if val.Name != "" {
				v.Outputs[i].Name = &val.Name
			}
			if val.Kind != "" {
				v.Outputs[i].Kind = &val.Kind
			}
			if val.Runnable != nil {
				v.Outputs[i].Runnable = protobufRunnablepbRunnableRefToRunnableRunnableRef(val.Runnable)
			}
			if val.Metadata != nil {
				v.Outputs[i].Metadata = protobufRunnablepbInputParameterToRunnableInputParameter(val.Metadata)
			}
		}
	}
	if message.Labels != nil {
		v.Labels = make([]string, len(message.Labels))
		for i, val := range message.Labels {
			v.Labels[i] = val
		}
	}
	return v
}

// NewRegisterResponse builds the gRPC response type from the result of the
// "register" endpoint of the "runnable" service.
func NewRegisterResponse(result *runnable.Runnable) *runnablepb.RegisterResponse {
	message := &runnablepb.RegisterResponse{
		Name: result.Name,
		Kind: result.Kind,
	}
	if result.ID != nil {
		message.Id = *result.ID
	}
	if result.Created != nil {
		message.Created = *result.Created
	}
	if result.Image != nil {
		message.Image = svcRunnableRunnableImageToRunnablepbRunnableImage(result.Image)
	}
	if result.Inputs != nil {
		message.Inputs = make([]*runnablepb.RunnableInput, len(result.Inputs))
		for i, val := range result.Inputs {
			message.Inputs[i] = &runnablepb.RunnableInput{}
			if val.Name != nil {
				message.Inputs[i].Name = *val.Name
			}
			if val.Kind != nil {
				message.Inputs[i].Kind = *val.Kind
			}
			if val.Runnable != nil {
				message.Inputs[i].Runnable = *val.Runnable
			}
		}
	}
	if result.Outputs != nil {
		message.Outputs = make([]*runnablepb.RunnableOutput, len(result.Outputs))
		for i, val := range result.Outputs {
			message.Outputs[i] = &runnablepb.RunnableOutput{}
			if val.Name != nil {
				message.Outputs[i].Name = *val.Name
			}
			if val.Kind != nil {
				message.Outputs[i].Kind = *val.Kind
			}
			if val.Runnable != nil {
				message.Outputs[i].Runnable = svcRunnableRunnableRefToRunnablepbRunnableRef(val.Runnable)
			}
			if val.Metadata != nil {
				message.Outputs[i].Metadata = svcRunnableInputParameterToRunnablepbInputParameter(val.Metadata)
			}
		}
	}
	if result.Labels != nil {
		message.Labels = make([]string, len(result.Labels))
		for i, val := range result.Labels {
			message.Labels[i] = val
		}
	}
	return message
}

// NewGetPayload builds the payload of the "get" endpoint of the "runnable"
// service from the gRPC request type.
func NewGetPayload(message *runnablepb.GetRequest) *runnable.GetPayload {
	v := &runnable.GetPayload{
		RunnableNameOrID: message.RunnableNameOrId,
	}
	return v
}

// NewGetResponse builds the gRPC response type from the result of the "get"
// endpoint of the "runnable" service.
func NewGetResponse(result *runnable.Runnable) *runnablepb.GetResponse {
	message := &runnablepb.GetResponse{
		Name: result.Name,
		Kind: result.Kind,
	}
	if result.ID != nil {
		message.Id = *result.ID
	}
	if result.Created != nil {
		message.Created = *result.Created
	}
	if result.Image != nil {
		message.Image = svcRunnableRunnableImageToRunnablepbRunnableImage(result.Image)
	}
	if result.Inputs != nil {
		message.Inputs = make([]*runnablepb.RunnableInput, len(result.Inputs))
		for i, val := range result.Inputs {
			message.Inputs[i] = &runnablepb.RunnableInput{}
			if val.Name != nil {
				message.Inputs[i].Name = *val.Name
			}
			if val.Kind != nil {
				message.Inputs[i].Kind = *val.Kind
			}
			if val.Runnable != nil {
				message.Inputs[i].Runnable = *val.Runnable
			}
		}
	}
	if result.Outputs != nil {
		message.Outputs = make([]*runnablepb.RunnableOutput, len(result.Outputs))
		for i, val := range result.Outputs {
			message.Outputs[i] = &runnablepb.RunnableOutput{}
			if val.Name != nil {
				message.Outputs[i].Name = *val.Name
			}
			if val.Kind != nil {
				message.Outputs[i].Kind = *val.Kind
			}
			if val.Runnable != nil {
				message.Outputs[i].Runnable = svcRunnableRunnableRefToRunnablepbRunnableRef(val.Runnable)
			}
			if val.Metadata != nil {
				message.Outputs[i].Metadata = svcRunnableInputParameterToRunnablepbInputParameter(val.Metadata)
			}
		}
	}
	if result.Labels != nil {
		message.Labels = make([]string, len(result.Labels))
		for i, val := range result.Labels {
			message.Labels[i] = val
		}
	}
	return message
}

// ValidateRunnableImage runs the validations defined on RunnableImage.
func ValidateRunnableImage(message *runnablepb.RunnableImage) (err error) {

	return
}

// ValidateRunnableInput runs the validations defined on RunnableInput.
func ValidateRunnableInput(message *runnablepb.RunnableInput) (err error) {

	return
}

// ValidateRunnableOutput runs the validations defined on RunnableOutput.
func ValidateRunnableOutput(message *runnablepb.RunnableOutput) (err error) {

	return
}

// ValidateRunnableRef runs the validations defined on RunnableRef.
func ValidateRunnableRef(message *runnablepb.RunnableRef) (err error) {

	return
}

// ValidateInputParameter runs the validations defined on InputParameter.
func ValidateInputParameter(message *runnablepb.InputParameter) (err error) {

	return
}

// ValidateRegisterRequest runs the validations defined on RegisterRequest.
func ValidateRegisterRequest(message *runnablepb.RegisterRequest) (err error) {
	if message.Image == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("image", "message"))
	}
	if message.Inputs == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("inputs", "message"))
	}
	if message.Outputs == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("outputs", "message"))
	}
	if message.Labels == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("labels", "message"))
	}
	if message.Id != "" {
		err = goa.MergeErrors(err, goa.ValidatePattern("message.id", message.Id, "uuid"))
	}
	if message.Created != "" {
		err = goa.MergeErrors(err, goa.ValidatePattern("message.created", message.Created, "date-time"))
	}
	return
}

// svcRunnableRunnableImageToRunnablepbRunnableImage builds a value of type
// *runnablepb.RunnableImage from a value of type *runnable.RunnableImage.
func svcRunnableRunnableImageToRunnablepbRunnableImage(v *runnable.RunnableImage) *runnablepb.RunnableImage {
	res := &runnablepb.RunnableImage{}
	if v.RegistryURL != nil {
		res.RegistryUrl = *v.RegistryURL
	}
	if v.Repository != nil {
		res.Repository = *v.Repository
	}
	if v.Tag != nil {
		res.Tag = *v.Tag
	}

	return res
}

// svcRunnableRunnableRefToRunnablepbRunnableRef builds a value of type
// *runnablepb.RunnableRef from a value of type *runnable.RunnableRef.
func svcRunnableRunnableRefToRunnablepbRunnableRef(v *runnable.RunnableRef) *runnablepb.RunnableRef {
	if v == nil {
		return nil
	}
	res := &runnablepb.RunnableRef{}
	if v.Name != nil {
		res.Name = *v.Name
	}
	if v.Kind != nil {
		res.Kind = *v.Kind
	}
	if v.Labels != nil {
		res.Labels = make([]string, len(v.Labels))
		for i, val := range v.Labels {
			res.Labels[i] = val
		}
	}

	return res
}

// svcRunnableInputParameterToRunnablepbInputParameter builds a value of type
// *runnablepb.InputParameter from a value of type *runnable.InputParameter.
func svcRunnableInputParameterToRunnablepbInputParameter(v *runnable.InputParameter) *runnablepb.InputParameter {
	if v == nil {
		return nil
	}
	res := &runnablepb.InputParameter{}
	if v.Datatype != nil {
		res.Datatype = *v.Datatype
	}
	if v.Optional != nil {
		res.Optional = *v.Optional
	}
	if v.Default != nil {
		res.Default = *v.Default
	}

	return res
}

// protobufRunnablepbRunnableImageToRunnableRunnableImage builds a value of
// type *runnable.RunnableImage from a value of type *runnablepb.RunnableImage.
func protobufRunnablepbRunnableImageToRunnableRunnableImage(v *runnablepb.RunnableImage) *runnable.RunnableImage {
	res := &runnable.RunnableImage{}
	if v.RegistryUrl != "" {
		res.RegistryURL = &v.RegistryUrl
	}
	if v.Repository != "" {
		res.Repository = &v.Repository
	}
	if v.Tag != "" {
		res.Tag = &v.Tag
	}

	return res
}

// protobufRunnablepbRunnableRefToRunnableRunnableRef builds a value of type
// *runnable.RunnableRef from a value of type *runnablepb.RunnableRef.
func protobufRunnablepbRunnableRefToRunnableRunnableRef(v *runnablepb.RunnableRef) *runnable.RunnableRef {
	if v == nil {
		return nil
	}
	res := &runnable.RunnableRef{}
	if v.Name != "" {
		res.Name = &v.Name
	}
	if v.Kind != "" {
		res.Kind = &v.Kind
	}
	if v.Labels != nil {
		res.Labels = make([]string, len(v.Labels))
		for i, val := range v.Labels {
			res.Labels[i] = val
		}
	}

	return res
}

// protobufRunnablepbInputParameterToRunnableInputParameter builds a value of
// type *runnable.InputParameter from a value of type
// *runnablepb.InputParameter.
func protobufRunnablepbInputParameterToRunnableInputParameter(v *runnablepb.InputParameter) *runnable.InputParameter {
	if v == nil {
		return nil
	}
	res := &runnable.InputParameter{
		Optional: &v.Optional,
	}
	if v.Datatype != "" {
		res.Datatype = &v.Datatype
	}
	if v.Default != "" {
		res.Default = &v.Default
	}

	return res
}
