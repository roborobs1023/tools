package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/roborobs1023/tools/internal/utils"
)

var (
	verifier             = emailverifier.NewVerifier().EnableAutoUpdateDisposable()
	domainRegex          = regexp.MustCompile(`(?i)^(?:([a-z0-9-]+|\*)\.)?([a-z0-9-]{1,61})\.([a-z0-9]{2,7})$`)
	nonNumericStartRegex = regexp.MustCompile(`^[A-Za-z]*[A-Za-z][A-Za-z0-9-. _]*$`)
)
var (
	disposableDomains = []string{"example.com", "example.org", "example.co", "example.net", "test.com", "test.org"}
)

type Config struct {
	DisableDisposableEmailCheck bool
	DisableCatchAllCheck        bool
}

func init() {
	verifier.AddDisposableDomains(disposableDomains)
}

func Validate(val interface{}, cfg *Config) error {
	var errs utils.Errs

	if cfg.DisableDisposableEmailCheck {
		verifier = verifier.DisableAutoUpdateDisposable()
	}

	if cfg.DisableCatchAllCheck {
		verifier = verifier.DisableCatchAllCheck()
	}

	v := reflect.ValueOf(val)
fieldsLoop:
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := v.Type().Field(i).Tag.Get("validate")

		if tag == "" {
			continue
		}
		fieldName := v.Type().Field(i).Name

		rules := strings.Split(tag, ",")

		for _, rule := range rules {
			switch {
			case rule == "required":
				if err := validateRequired(field, fieldName); err != nil {
					errs.Add(err)
					continue fieldsLoop
				}
			case strings.HasPrefix(rule, "min="):
				if err := validateMinLength(rule, field, fieldName); err != nil {
					errs.Add(err)
					continue fieldsLoop
				}
			case strings.HasPrefix(rule, "max="):
				if err := validateMaxLength(rule, field, fieldName); err != nil {
					errs.Add(err)
				}

			case rule == "email":
				if err := validateEmail(field, fieldName); err != nil {
					errs.Add(err)
					continue fieldsLoop
				}

			case strings.HasPrefix(rule, "req_domain="):
				if err := validateRequiredDomain(rule, field, fieldName); err != nil {
					errs.Add(err)
				}

			case rule == "nonDisposable":
				if err := validateNonDisposableDomain(rules, field, fieldName); err != nil {
					errs.Add(err)
				}
			case rule == "domain":
				if err := validateDomain(field, fieldName); err != nil {
					errs.Add(err)
				}

			case rule == "nonNumericStart":
				if err := validateNonNumericStart(field, fieldName); err != nil {
					errs.Add(err)
				}
			}
		}
	}
	if errs.Len() != 0 {
		return errs
	}

	return nil
}

func validateNonNumericStart(field reflect.Value, fieldName string) error {
	if !nonNumericStartRegex.MatchString(field.String()) {
		return fmt.Errorf("%s must start with a letter", fieldName)
	}
	return nil
}

func validateRequired(field reflect.Value, fieldName string) error {
	if field.String() == "" {
		return fmt.Errorf(
			"%s is a required field",
			fieldName,
		)
	}
	return nil
}

func validateMinLength(rule string, field reflect.Value, fieldName string) error {
	min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
	if len(field.String()) < min {
		return fmt.Errorf(
			"%s must be at least %d characters",
			fieldName,
			min,
		)
	}
	return nil
}

func validateMaxLength(rule string, field reflect.Value, fieldName string) error {
	max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
	if len(field.String()) > max {
		return fmt.Errorf(
			"%s must be less than %d characters",
			fieldName,
			max,
		)
	}
	return nil
}

func validateEmail(field reflect.Value, fieldName string) error {
	ret, err := verifier.Verify(field.String())
	if err != nil {
		return fmt.Errorf(
			"%s must be a valid email address",
			fieldName,
		)
	}

	if !ret.Syntax.Valid {
		return fmt.Errorf("invalid email syntax")
	}
	return nil
}

func validateRequiredDomain(rule string, field reflect.Value, fieldName string) error {
	ret, err := verifier.Verify(field.String())

	if err != nil {
		return fmt.Errorf(
			"%s must be a valid email address",
			fieldName,
		)
	}

	domain := strings.TrimPrefix(rule, "req_domain=")

	if ret.Syntax.Domain != domain {
		return fmt.Errorf(
			"%s must be a part of the '%s' domain",
			fieldName,
			domain,
		)
	}

	return nil
}

func validateNonDisposableDomain(rules []string, field reflect.Value, fieldName string) error {
	var domain = ""

	if slices.Contains(rules, "email") {
		ret, err := verifier.Verify(field.String())

		if err != nil {
			return fmt.Errorf(
				"%s must be a valid email address",
				fieldName,
			)
		}
		domain = ret.Syntax.Domain
	} else {
		domain = field.String()
	}

	if verifier.IsDisposable(domain) {
		return fmt.Errorf("%s contains a disposable domain", fieldName)
	}

	return nil
}

func validateDomain(field reflect.Value, fieldName string) error {
	if !domainRegex.MatchString(field.String()) {
		return fmt.Errorf("%s is an invalid domain", fieldName)
	}

	return nil
}
