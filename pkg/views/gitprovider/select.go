// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package gitprovider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/daytonaio/daytona/cmd/daytona/config"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"
)

var commonGitProviderIds = []string{"github", "gitlab", "bitbucket"}

func GitProviderSelectionView(ctx context.Context, gitProviderAddView *apiclient.SetGitProviderConfig, apiClient *apiclient.APIClient) {
	supportedProviders := config.GetSupportedGitProviders()

	var gitProviderOptions []huh.Option[string]
	var otherGitProviderOptions []huh.Option[string]
	for _, supportedProvider := range supportedProviders {
		if slices.Contains(commonGitProviderIds, supportedProvider.Id) {
			gitProviderOptions = append(gitProviderOptions, huh.Option[string]{Key: supportedProvider.Name, Value: supportedProvider.Id})
		} else {
			otherGitProviderOptions = append(otherGitProviderOptions, huh.Option[string]{Key: supportedProvider.Name, Value: supportedProvider.Id})
		}
	}

	if len(otherGitProviderOptions) > 0 {
		gitProviderOptions = append(gitProviderOptions, huh.Option[string]{Key: "Other", Value: "other"})
	}

	gitProviderForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Git provider").
				Options(
					gitProviderOptions...,
				).
				Value(&gitProviderAddView.ProviderId)).WithHeight(8),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Git provider").
				Options(
					otherGitProviderOptions...,
				).
				Value(&gitProviderAddView.ProviderId)).WithHeight(12).WithHideFunc(func() bool {
			return gitProviderAddView.ProviderId != "other"
		}),
	).WithTheme(views.GetCustomTheme())

	err := gitProviderForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	var selectedSigningMethod string
	var signingKey string

	userDataForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Value(gitProviderAddView.Username).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("username can not be blank")
					}
					return nil
				}),
		).WithHeight(5).WithHideFunc(func() bool {
			return !providerRequiresUsername(gitProviderAddView.ProviderId)
		}),
		huh.NewGroup(
			huh.NewInput().
				Title("Self-managed API URL").
				Value(gitProviderAddView.BaseApiUrl).
				Description(getApiUrlDescription(gitProviderAddView.ProviderId)).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("URL can not be blank")
					}
					return nil
				}),
		).WithHeight(6).WithHideFunc(func() bool {
			return !providerRequiresApiUrl(gitProviderAddView.ProviderId)
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("Personal access token").
				Value(&gitProviderAddView.Token).
				EchoMode(huh.EchoModePassword).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("token can not be blank")
					}
					return nil
				}),
		).WithHeight(5),
		huh.NewGroup(
			huh.NewInput().
				Title("Alias").
				Description("Will default to username if left empty").
				Value(gitProviderAddView.Alias),
		).WithHeight(6),

		huh.NewGroup(huh.NewSelect[string]().
			Title("Commit Signing Method").
			DescriptionFunc(func() string {
				return getGitProviderSigningHelpMessage(*gitProviderAddView.Id)
			}, nil).
			Options(
				huh.Option[string]{Key: "None", Value: "none"},
				huh.Option[string]{Key: "SSH", Value: "ssh"},
				huh.Option[string]{Key: "GPG", Value: "gpg"},
			).
			Value(&selectedSigningMethod).WithHeight(6),
		).WithHeight(8).WithHideFunc(func() bool {
			return commitSigningNotSupported(*gitProviderAddView.Id)
		}),
		huh.NewGroup(
			huh.NewInput().
				Title("Signing Key").
				Value(&signingKey).
				DescriptionFunc(func() string {
					return getSigningKeyDescription(selectedSigningMethod)
				}, nil).
				Validate(func(str string) error {
					if selectedSigningMethod != "none" && str == "" {
						return errors.New("signing key cannot be blank when a signing method is selected")
					}

					if selectedSigningMethod == "ssh" {
						if err := isValidSSHKey(str); err != nil {
							return err
						}
					}
					return nil
				}),
		).WithHeight(5).WithHideFunc(func() bool {
			return selectedSigningMethod == "none"
		}),
	).WithTheme(views.GetCustomTheme())

	views.RenderInfoMessage(getGitProviderHelpMessage(gitProviderAddView.ProviderId))
	err = userDataForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	if selectedSigningMethod != "none" {
		gitProviderAddView.SigningMethod = (*apiclient.SigningMethod)(&selectedSigningMethod)
		gitProviderAddView.SigningKey = &signingKey
	}
}
func isValidSSHKey(key string) error {
	sshKeyPattern := regexp.MustCompile(`^(ssh-(rsa|ed25519|dss|ecdsa-sha2-nistp(256|384|521)))\s+[A-Za-z0-9+/=]+(\s+.+)?$`)
	if !sshKeyPattern.MatchString(key) {
		return errors.New("invalid SSH key: must start with valid SSH key type (e.g., ssh-rsa, ssh-ed25519)")
	}
	return nil
}

func GitProviderDeleteView(gitProviderAddView *apiclient.SetGitProviderConfig, userGitProviders []apiclient.GitProvider, apiClient *apiclient.APIClient, ctx context.Context) {
	var gitProviderOptions []huh.Option[string]
	for _, userProvider := range userGitProviders {
		gitProviderOptions = append(gitProviderOptions, huh.Option[string]{Key: fmt.Sprintf("%s   %s", userProvider.ProviderId, userProvider.Alias), Value: userProvider.Id})
	}
	gitProviderForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Git provider").
				Options(
					gitProviderOptions...,
				).
				Value(gitProviderAddView.Id)).WithHeight(8),
	).WithTheme(views.GetCustomTheme())
	err := gitProviderForm.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func providerRequiresUsername(gitProviderId string) bool {
	return gitProviderId == "bitbucket" || gitProviderId == "bitbucket-server" || gitProviderId == "aws-codecommit"
}

func providerRequiresApiUrl(gitProviderId string) bool {
	return gitProviderId == "gitness" || gitProviderId == "github-enterprise-server" || gitProviderId == "gitlab-self-managed" || gitProviderId == "gitea" || gitProviderId == "bitbucket-server" || gitProviderId == "azure-devops" || gitProviderId == "aws-codecommit"
}

func commitSigningNotSupported(gitProviderId string) bool {
	return gitProviderId == "gitness" || gitProviderId == "bitbucket" || gitProviderId == "bitbucket-server"
}

func getApiUrlDescription(gitProviderId string) string {
	if gitProviderId == "gitlab-self-managed" {
		return "For example: http://gitlab-host/api/v4/"
	} else if gitProviderId == "github-enterprise-server" {
		return "For example: https://github-host"
	} else if gitProviderId == "gitea" {
		return "For example: http://gitea-host"
	} else if gitProviderId == "gitness" {
		return "For example: http://gitness-host/api/v1/"
	} else if gitProviderId == "azure-devops" {
		return "For example: https://dev.azure.com/organization"
	} else if gitProviderId == "bitbucket-server" {
		return "For example: https://bitbucket.host.com/rest"
	} else if gitProviderId == "aws-codecommit" {
		return "For example: https://ap-south-1.console.aws.amazon.com"
	}
	return ""
}

func getSigningKeyDescription(signingMethod string) string {
	switch signingMethod {
	case "gpg":
		return "Provide your GPG key ID (e.g., 30F2B65B9246B6CA) for signing commits."
	case "ssh":
		return "Provide your public SSH key (e.g., ssh-ed25519 AAAAC3...<rest of key>) for secure signing."
	default:
		return ""
	}
}

func getGitProviderHelpMessage(gitProviderId string) string {
	message := fmt.Sprintf("%s\n%s\n\n%s%s",
		lipgloss.NewStyle().Foreground(views.Green).Bold(true).Render("More information on:"),
		config.GetDocsLinkFromGitProvider(gitProviderId),
		lipgloss.NewStyle().Foreground(views.Green).Bold(true).Render("Required scopes: "),
		config.GetRequiredScopesFromGitProviderId(gitProviderId))

	prebuildScopes := config.GetPrebuildScopesFromGitProviderId(gitProviderId)
	if prebuildScopes != "" {
		message = fmt.Sprintf("%s\n%s%s",
			message,
			lipgloss.NewStyle().Foreground(views.Green).Bold(true).Render("Prebuild scopes: "),
			prebuildScopes)
	}

	return message
}

func getGitProviderSigningHelpMessage(gitProviderId string) string {
	signingDocsLink := config.GetDocsLinkForCommitSigning(gitProviderId)

	if signingDocsLink != "" {
		return signingDocsLink
	}
	return ""
}
