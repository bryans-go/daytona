/*
Daytona Server API

Daytona Server API

API version: v0.0.0-dev
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// checks if the CreateProjectDTO type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateProjectDTO{}

// CreateProjectDTO struct for CreateProjectDTO
type CreateProjectDTO struct {
	BuildConfig         *BuildConfig           `json:"buildConfig,omitempty"`
	EnvVars             map[string]string      `json:"envVars"`
	GitProviderConfigId *string                `json:"gitProviderConfigId,omitempty"`
	Image               *string                `json:"image,omitempty"`
	Name                string                 `json:"name"`
	Source              CreateProjectSourceDTO `json:"source"`
	User                *string                `json:"user,omitempty"`
}

type _CreateProjectDTO CreateProjectDTO

// NewCreateProjectDTO instantiates a new CreateProjectDTO object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateProjectDTO(envVars map[string]string, name string, source CreateProjectSourceDTO) *CreateProjectDTO {
	this := CreateProjectDTO{}
	this.EnvVars = envVars
	this.Name = name
	this.Source = source
	return &this
}

// NewCreateProjectDTOWithDefaults instantiates a new CreateProjectDTO object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateProjectDTOWithDefaults() *CreateProjectDTO {
	this := CreateProjectDTO{}
	return &this
}

// GetBuildConfig returns the BuildConfig field value if set, zero value otherwise.
func (o *CreateProjectDTO) GetBuildConfig() BuildConfig {
	if o == nil || IsNil(o.BuildConfig) {
		var ret BuildConfig
		return ret
	}
	return *o.BuildConfig
}

// GetBuildConfigOk returns a tuple with the BuildConfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateProjectDTO) GetBuildConfigOk() (*BuildConfig, bool) {
	if o == nil || IsNil(o.BuildConfig) {
		return nil, false
	}
	return o.BuildConfig, true
}

// HasBuildConfig returns a boolean if a field has been set.
func (o *CreateProjectDTO) HasBuildConfig() bool {
	if o != nil && !IsNil(o.BuildConfig) {
		return true
	}

	return false
}

// SetBuildConfig gets a reference to the given BuildConfig and assigns it to the BuildConfig field.
func (o *CreateProjectDTO) SetBuildConfig(v BuildConfig) {
	o.BuildConfig = &v
}

// GetEnvVars returns the EnvVars field value
func (o *CreateProjectDTO) GetEnvVars() map[string]string {
	if o == nil {
		var ret map[string]string
		return ret
	}

	return o.EnvVars
}

// GetEnvVarsOk returns a tuple with the EnvVars field value
// and a boolean to check if the value has been set.
func (o *CreateProjectDTO) GetEnvVarsOk() (*map[string]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EnvVars, true
}

// SetEnvVars sets field value
func (o *CreateProjectDTO) SetEnvVars(v map[string]string) {
	o.EnvVars = v
}

// GetGitProviderConfigId returns the GitProviderConfigId field value if set, zero value otherwise.
func (o *CreateProjectDTO) GetGitProviderConfigId() string {
	if o == nil || IsNil(o.GitProviderConfigId) {
		var ret string
		return ret
	}
	return *o.GitProviderConfigId
}

// GetGitProviderConfigIdOk returns a tuple with the GitProviderConfigId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateProjectDTO) GetGitProviderConfigIdOk() (*string, bool) {
	if o == nil || IsNil(o.GitProviderConfigId) {
		return nil, false
	}
	return o.GitProviderConfigId, true
}

// HasGitProviderConfigId returns a boolean if a field has been set.
func (o *CreateProjectDTO) HasGitProviderConfigId() bool {
	if o != nil && !IsNil(o.GitProviderConfigId) {
		return true
	}

	return false
}

// SetGitProviderConfigId gets a reference to the given string and assigns it to the GitProviderConfigId field.
func (o *CreateProjectDTO) SetGitProviderConfigId(v string) {
	o.GitProviderConfigId = &v
}

// GetImage returns the Image field value if set, zero value otherwise.
func (o *CreateProjectDTO) GetImage() string {
	if o == nil || IsNil(o.Image) {
		var ret string
		return ret
	}
	return *o.Image
}

// GetImageOk returns a tuple with the Image field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateProjectDTO) GetImageOk() (*string, bool) {
	if o == nil || IsNil(o.Image) {
		return nil, false
	}
	return o.Image, true
}

// HasImage returns a boolean if a field has been set.
func (o *CreateProjectDTO) HasImage() bool {
	if o != nil && !IsNil(o.Image) {
		return true
	}

	return false
}

// SetImage gets a reference to the given string and assigns it to the Image field.
func (o *CreateProjectDTO) SetImage(v string) {
	o.Image = &v
}

// GetName returns the Name field value
func (o *CreateProjectDTO) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *CreateProjectDTO) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *CreateProjectDTO) SetName(v string) {
	o.Name = v
}

// GetSource returns the Source field value
func (o *CreateProjectDTO) GetSource() CreateProjectSourceDTO {
	if o == nil {
		var ret CreateProjectSourceDTO
		return ret
	}

	return o.Source
}

// GetSourceOk returns a tuple with the Source field value
// and a boolean to check if the value has been set.
func (o *CreateProjectDTO) GetSourceOk() (*CreateProjectSourceDTO, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Source, true
}

// SetSource sets field value
func (o *CreateProjectDTO) SetSource(v CreateProjectSourceDTO) {
	o.Source = v
}

// GetUser returns the User field value if set, zero value otherwise.
func (o *CreateProjectDTO) GetUser() string {
	if o == nil || IsNil(o.User) {
		var ret string
		return ret
	}
	return *o.User
}

// GetUserOk returns a tuple with the User field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateProjectDTO) GetUserOk() (*string, bool) {
	if o == nil || IsNil(o.User) {
		return nil, false
	}
	return o.User, true
}

// HasUser returns a boolean if a field has been set.
func (o *CreateProjectDTO) HasUser() bool {
	if o != nil && !IsNil(o.User) {
		return true
	}

	return false
}

// SetUser gets a reference to the given string and assigns it to the User field.
func (o *CreateProjectDTO) SetUser(v string) {
	o.User = &v
}

func (o CreateProjectDTO) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateProjectDTO) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.BuildConfig) {
		toSerialize["buildConfig"] = o.BuildConfig
	}
	toSerialize["envVars"] = o.EnvVars
	if !IsNil(o.GitProviderConfigId) {
		toSerialize["gitProviderConfigId"] = o.GitProviderConfigId
	}
	if !IsNil(o.Image) {
		toSerialize["image"] = o.Image
	}
	toSerialize["name"] = o.Name
	toSerialize["source"] = o.Source
	if !IsNil(o.User) {
		toSerialize["user"] = o.User
	}
	return toSerialize, nil
}

func (o *CreateProjectDTO) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"envVars",
		"name",
		"source",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varCreateProjectDTO := _CreateProjectDTO{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varCreateProjectDTO)

	if err != nil {
		return err
	}

	*o = CreateProjectDTO(varCreateProjectDTO)

	return err
}

type NullableCreateProjectDTO struct {
	value *CreateProjectDTO
	isSet bool
}

func (v NullableCreateProjectDTO) Get() *CreateProjectDTO {
	return v.value
}

func (v *NullableCreateProjectDTO) Set(val *CreateProjectDTO) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateProjectDTO) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateProjectDTO) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateProjectDTO(val *CreateProjectDTO) *NullableCreateProjectDTO {
	return &NullableCreateProjectDTO{value: val, isSet: true}
}

func (v NullableCreateProjectDTO) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateProjectDTO) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
