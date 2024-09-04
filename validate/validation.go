package validate

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/roborobs1023/tools/internal/utils"
)

var (
	verifier = emailverifier.NewVerifier().EnableAutoUpdateDisposable()
	// domainRegex          = regexp.MustCompile(`(?i)^(?:([a-z0-9-]+|\*)\.)?([a-z0-9-]{1,61})\.([a-z0-9]{2,7})$`)
	nonNumericStartRegex = regexp.MustCompile(`^[A-Za-z]*[A-Za-z][A-Za-z0-9-. _]*$`)
)
var (
	disposableDomains = []string{"example.com", "example.org", "example.co", "example.net", "test.com", "test.org"}
)

type Config struct {
	DisableDisposableEmailCheck bool
	DisableCatchAllCheck        bool
	IgnoreEmptyFields           bool
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

	if !cfg.IgnoreEmptyFields {
		cfg.IgnoreEmptyFields = true
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

		if field.IsZero() && slices.Contains(rules, "optional") {
			continue
		}
		for _, rule := range rules {

			switch {
			case rule == "required":
				if err := validateRequired(field, fieldName); err != nil {
					errs.Add(err)
					continue fieldsLoop
				}
			case strings.HasPrefix(rule, "min="):
				if err := validateMin(rule, field, fieldName); err != nil {
					errs.Add(err)
					continue fieldsLoop
				}
			case strings.HasPrefix(rule, "max="):
				if err := validateMax(rule, field, fieldName); err != nil {
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

func validateMin(rule string, field reflect.Value, fieldName string) error {
	if field.CanInt() {
		min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
		i := field.Int()

		if i < int64(min) {

			return fmt.Errorf("%s must be more than %v", fieldName, min)
		} else {
			return nil
		}
	} else if field.CanFloat() {
		min, _ := strconv.ParseFloat(strings.TrimPrefix(rule, "min="), 64)
		f := field.Float()

		if f < min {
			return fmt.Errorf("%s must be more than %v", fieldName, min)
		}
	} else {
		min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
		if len(field.String()) < min {
			return fmt.Errorf(
				"%s must be more than %d characters",
				fieldName,
				min,
			)
		}
	}
	return nil
}

func validateMax(rule string, field reflect.Value, fieldName string) error {

	if field.CanInt() {
		max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
		i := field.Int()

		if i > int64(max) {
			return fmt.Errorf("%s must be less than %v", fieldName, max)
		}
	} else if field.CanFloat() {
		max, _ := strconv.ParseFloat(strings.TrimPrefix(rule, "max="), 64)
		f := field.Float()

		if f > max {
			return fmt.Errorf("%s must be less than %v", fieldName, max)
		}
	} else {
		max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
		if len(field.String()) > max {
			return fmt.Errorf(
				"%s must be less than %d characters",
				fieldName,
				max,
			)
		}
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

	// err = validateNonDisposableDomain([]string{"email"}, field, fieldName)
	// if err != nil {
	// 	return fmt.Errorf("email has an invalid domain")
	// }

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

	log.Println(domain)
	if verifier.IsDisposable(domain) {
		return fmt.Errorf("%s contains a disposable domain", fieldName)
	}

	return nil
}

func validateDomain(field reflect.Value, fieldName string) error {
	exist, err := CheckDomain(field.String())
	if !exist {
		log.Printf("%s error: %s", fieldName, err)
		return fmt.Errorf("%s is an invalid domain", fieldName)
	}

	return nil
}
