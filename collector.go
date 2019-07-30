package main

import (
	"bareos_exporter/DataAccess/Queries"
	"bareos_exporter/Error"
	"github.com/prometheus/client_golang/prometheus"
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type BareosMetrics struct {
	TotalFiles *prometheus.Desc
	TotalBytes *prometheus.Desc
	LastJobBytes *prometheus.Desc
	LastJobFiles *prometheus.Desc
	LastJobErrors *prometheus.Desc
	LastJobTimestamp *prometheus.Desc

	LastFullJobBytes *prometheus.Desc
	LastFullJobFiles *prometheus.Desc
	LastFullJobErrors *prometheus.Desc
	LastFullJobTimestamp *prometheus.Desc
}

func BareosCollector() *BareosMetrics {
	return &BareosMetrics{
		TotalFiles: prometheus.NewDesc("total_files_saved_for_hostname",
			"Total files saved for server during all backups combined",
			[]string{"hostname"}, nil,
		),
		TotalBytes: prometheus.NewDesc("total_bytes_saved_for_hostname",
			"Total bytes saved for server during all backups combined",
			[]string{"hostname"}, nil,
		),
		LastJobBytes: prometheus.NewDesc("last_backup_job_bytes_saved_for_hostname",
			"Last job that was executed for ",
			[]string{"hostname"}, nil,
		),
		LastJobFiles: prometheus.NewDesc("last_backup_job_files_saved_for_hostname",
			"Last job that was executed for ",
			[]string{"hostname"}, nil,
		),
		LastJobErrors: prometheus.NewDesc("last_backup_job_errors_occurred_while_saving_for_hostname",
			"Last job that was executed for ",
			[]string{"hostname"}, nil,
		),
		LastJobTimestamp: prometheus.NewDesc("last_backup_job_unix_timestamp_for_hostname",
			"Last job that was executed for ",
			[]string{"hostname"}, nil,
		),
		LastFullJobBytes: prometheus.NewDesc("last_full_backup_job_bytes_saved_for_hostname",
			"Total bytes saved by server",
			[]string{"hostname"}, nil,
		),
		LastFullJobFiles: prometheus.NewDesc("last_full_backup_job_files_saved_for_hostname",
			"Total bytes saved by server",
			[]string{"hostname"}, nil,
		),
		LastFullJobErrors: prometheus.NewDesc("last_full_backup_job_errors_occurred_while_saving_for_hostname",
			"Total bytes saved by server",
			[]string{"hostname"}, nil,
		),
		LastFullJobTimestamp: prometheus.NewDesc("last_full_backup_job_unix_timestamp_for_hostname",
			"Total bytes saved by server",
			[]string{"hostname"}, nil,
		),
	}
}

func (collector *BareosMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.TotalFiles
	ch <- collector.TotalBytes
	ch <- collector.LastJobBytes
	ch <- collector.LastJobFiles
	ch <- collector.LastJobErrors
	ch <- collector.LastJobTimestamp
	ch <- collector.LastFullJobBytes
	ch <- collector.LastFullJobFiles
	ch <- collector.LastFullJobErrors
	ch <- collector.LastFullJobTimestamp
}

func (collector *BareosMetrics) Collect(ch chan<- prometheus.Metric) {
	results, err := db.Query("SELECT DISTINCT Name FROM job WHERE SchedTime LIKE '2019-07-24%'")

	Error.Check(err)

	var servers []Queries.Server

	for results.Next() {
		var server Queries.Server
		err = results.Scan(&server.Name)

		servers = append(servers, server)

		Error.Check(err)
	}

	for _, server := range servers {
		serverFiles := server.TotalFiles(db)
		serverBytes := server.TotalBytes(db)
		lastServerJob := server.LastJob(db, false)

		ch <- prometheus.MustNewConstMetric(collector.TotalFiles, prometheus.CounterValue, float64(serverFiles.Files), server.Name)
		ch <- prometheus.MustNewConstMetric(collector.TotalBytes, prometheus.CounterValue, float64(serverBytes.Bytes), server.Name)

		ch <- prometheus.MustNewConstMetric(collector.LastJobBytes, prometheus.CounterValue, float64(lastServerJob.JobBytes), server.Name)
		ch <- prometheus.MustNewConstMetric(collector.LastJobFiles, prometheus.CounterValue, float64(lastServerJob.JobFiles), server.Name)
		ch <- prometheus.MustNewConstMetric(collector.LastJobErrors, prometheus.CounterValue, float64(lastServerJob.JobErrors), server.Name)
		ch <- prometheus.MustNewConstMetric(collector.LastJobTimestamp, prometheus.CounterValue, float64(lastServerJob.JobDate.Unix()), server.Name)
	}

	for _, server := range servers {
		lastFullServerJob := server.LastJob(db, true)

		ch <- prometheus.MustNewConstMetric(collector.LastFullJobBytes, prometheus.CounterValue, float64(lastFullServerJob.JobBytes), server.Name)
		ch <- prometheus.MustNewConstMetric(collector.LastFullJobFiles, prometheus.CounterValue, float64(lastFullServerJob.JobFiles), server.Name)
		ch <- prometheus.MustNewConstMetric(collector.LastFullJobErrors, prometheus.CounterValue, float64(lastFullServerJob.JobErrors), server.Name)
		ch <- prometheus.MustNewConstMetric(collector.LastFullJobTimestamp, prometheus.CounterValue, float64(lastFullServerJob.JobDate.Unix()), server.Name)
	}
}