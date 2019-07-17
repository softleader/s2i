package jenkins

import (
	"fmt"
)

const (
	pathJobBuild               = "/job/%s/build"
	pathJobBuildWithParameters = "/job/%s/buildWithParameters"
)

var (
	// ErrJobNotFound 表示 job name 不存在於 jenkins 上
	ErrJobNotFound = fmt.Errorf("job not found")
)

// Build build 傳入的 job
func (jc *JobClient) Build(jobName string) error {
	jc.c.log.Printf("enqueuing job %q to build queue", jobName)

	r, err := jc.c.csrf()
	if err != nil {
		return err
	}
	resp, err := r.Post(fmt.Sprintf(pathJobBuild, jobName))
	if err != nil {
		return err
	}
	if resp.StatusCode() == 404 {
		return ErrJobNotFound
	}
	if !resp.IsSuccess() {
		return ErrNot2xxStatusCode
	}
	if body := resp.Body(); len(body) > 0 {
		jc.c.log.Println(body)
	}
	return nil
}

// BuildWithParameters build 傳入的 job 及 parameters
func (jc *JobClient) BuildWithParameters(jobName string, params map[string]string) error {
	jc.c.log.Printf("enqueuing job %q to build queue", jobName)

	r, err := jc.c.csrf()
	if err != nil {
		return err
	}
	resp, err := r.
		SetQueryParams(params).
		Post(fmt.Sprintf(pathJobBuildWithParameters, jobName))
	if err != nil {
		return err
	}
	if resp.StatusCode() == 404 {
		return ErrJobNotFound
	}
	if !resp.IsSuccess() {
		return ErrNot2xxStatusCode
	}
	if body := resp.Body(); len(body) > 0 {
		jc.c.log.Println(body)
	}
	return nil
}
