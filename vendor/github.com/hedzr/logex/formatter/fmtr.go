/*
 * Copyright © 2019 Hedzr Yeh.
 */

package formatter

import "github.com/sirupsen/logrus"

func resolve(fieldMap logrus.FieldMap, key string) string {
	// if k, ok := fieldMap[(logrus.fieldKey)key]; ok {
	// 	return k
	// }

	return string(key)
}

func prefixFieldClashes(data logrus.Fields, fieldMap logrus.FieldMap, reportCaller bool) {
	timeKey := resolve(fieldMap, logrus.FieldKeyTime)
	if t, ok := data[timeKey]; ok {
		data["fields."+timeKey] = t
		delete(data, timeKey)
	}

	msgKey := resolve(fieldMap, logrus.FieldKeyMsg)
	if m, ok := data[msgKey]; ok {
		data["fields."+msgKey] = m
		delete(data, msgKey)
	}

	levelKey := resolve(fieldMap, logrus.FieldKeyLevel)
	if l, ok := data[levelKey]; ok {
		data["fields."+levelKey] = l
		delete(data, levelKey)
	}

	logrusErrKey := resolve(fieldMap, logrus.FieldKeyLogrusError)
	if l, ok := data[logrusErrKey]; ok {
		data["fields."+logrusErrKey] = l
		delete(data, logrusErrKey)
	}

	// If reportCaller is not set, 'func' will not conflict.
	if reportCaller {
		funcKey := resolve(fieldMap, logrus.FieldKeyFunc)
		if l, ok := data[funcKey]; ok {
			data["fields."+funcKey] = l
		}
		fileKey := resolve(fieldMap, logrus.FieldKeyFile)
		if l, ok := data[fileKey]; ok {
			data["fields."+fileKey] = l
		}
	}
}
