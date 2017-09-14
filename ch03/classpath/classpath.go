package classpath

import (
	"os"
	"path/filepath"
)

type Classpath struct {
	bootClasspath Entry
	extClasspath  Entry
	userClasspath Entry
}

func Parse(jreOption, cpOption string) *Classpath {
	cp := &Classpath{}
	cp.parseBootAndExtClassPath(jreOption)
	cp.parseUserClasspath(cpOption)
	return cp
}

func (this *Classpath) ReadClass(classname string) ([]byte, Entry, error) {
	classname = classname + ".class"
	if data, entry, err := this.bootClasspath.readClass(classname); err == nil {
		return data, entry, err
	}
	if data, entry, err := this.extClasspath.readClass(classname); err == nil {
		return data, entry, err
	}
	return this.userClasspath.readClass(classname)
}

func (this *Classpath) String() string {
	return this.userClasspath.String()
}
func (this *Classpath) parseBootAndExtClassPath(jreOption string) {
	jreDir := this.getJreDir(jreOption)
	jreLibPath := filepath.Join(jreDir, "lib", "*")
	this.bootClasspath = newWildcardEntry(jreLibPath)
	jreExtPath := filepath.Join(jreDir, "lib", "ext", "*")
	this.extClasspath = newWildcardEntry(jreExtPath)
}
func (this *Classpath) getJreDir(jreOption string) string {
	if jreOption != "" && exists(jreOption) {
		return jreOption
	}
	if exists("./jre") {
		return ".jre"
	}
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		return filepath.Join(jh, "jre")
	}
	panic("Can not find jre folder")
}

func (this *Classpath) parseUserClasspath(cpOption string) {
	if cpOption == "" {
		cpOption = "."
	}
	this.userClasspath = newEntry(cpOption)
}
func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
