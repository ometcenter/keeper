/* Package config определяет структуры конфигурации и некоторые интерфейсы
 */

package config

import (
	"testing"
	"time"
)

func TestServiceConfig_InitTimezone(t *testing.T) {
	time.Local, _ = time.LoadLocation("UTC")
	tests := []struct {
		name string
		s    ServiceConfig
		want string
	}{
		{
			name: "first case",
			s:    *New(),
			want: "[INFO] текущая таймзона UTC\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.InitTimezone(); got != tt.want {
				t.Errorf("ServiceConfig.InitTimezone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceConfig_InitTimezone_MSK(t *testing.T) {
	time.Local, _ = time.LoadLocation("Europe/Moscow")
	tests := []struct {
		name string
		s    ServiceConfig
		want string
	}{
		{
			name: "first case",
			s:    *New(),
			want: "[INFO] текущая таймзона MSK\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.InitTimezone(); got != tt.want {
				t.Errorf("ServiceConfig.InitTimezone() = %v, want %v", got, tt.want)
			}
		})
	}
}
