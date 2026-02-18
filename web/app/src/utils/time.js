import { translate } from '@/i18n'

/**
 * Generates a human-readable relative time string (e.g., "2 hours ago")
 * @param {string|Date} timestamp - The timestamp to convert
 * @returns {string} Relative time string
 */
export const generatePrettyTimeAgo = (timestamp) => {
  const differenceInMs = new Date().getTime() - new Date(timestamp).getTime()
  if (differenceInMs < 500) {
    return translate('time.now')
  }
  if (differenceInMs > 3 * 86400000) { // If it was more than 3 days ago, we'll display the number of days ago
    const days = Number((differenceInMs / 86400000).toFixed(0))
    return translate(days === 1 ? 'time.dayAgo' : 'time.daysAgo', { value: days })
  }
  if (differenceInMs > 3600000) { // If it was more than 1h ago, display the number of hours ago
    const hours = Number((differenceInMs / 3600000).toFixed(0))
    return translate(hours === 1 ? 'time.hourAgo' : 'time.hoursAgo', { value: hours })
  }
  if (differenceInMs > 60000) {
    const minutes = Number((differenceInMs / 60000).toFixed(0))
    return translate(minutes === 1 ? 'time.minuteAgo' : 'time.minutesAgo', { value: minutes })
  }
  const seconds = Number((differenceInMs / 1000).toFixed(0))
  return translate(seconds === 1 ? 'time.secondAgo' : 'time.secondsAgo', { value: seconds })
}

/**
 * Generates a pretty time difference string between two timestamps
 * @param {string|Date} start - Start timestamp
 * @param {string|Date} end - End timestamp
 * @returns {string} Time difference string
 */
export const generatePrettyTimeDifference = (start, end) => {
  const ms = new Date(start) - new Date(end)
  const seconds = Math.floor(ms / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)

  if (hours > 0) {
    const remainingMinutes = minutes % 60
    const hoursText = translate(hours === 1 ? 'time.hour' : 'time.hours', { value: hours })
    if (remainingMinutes > 0) {
      return `${hoursText} ${translate(remainingMinutes === 1 ? 'time.minute' : 'time.minutes', { value: remainingMinutes })}`
    }
    return hoursText
  } else if (minutes > 0) {
    const remainingSeconds = seconds % 60
    const minutesText = translate(minutes === 1 ? 'time.minute' : 'time.minutes', { value: minutes })
    if (remainingSeconds > 0) {
      return `${minutesText} ${translate(remainingSeconds === 1 ? 'time.second' : 'time.seconds', { value: remainingSeconds })}`
    }
    return minutesText
  } else {
    return translate(seconds === 1 ? 'time.second' : 'time.seconds', { value: seconds })
  }
}

/**
 * Formats a timestamp into YYYY-MM-DD HH:mm:ss format
 * @param {string|Date} timestamp - The timestamp to format
 * @returns {string} Formatted timestamp
 */
export const prettifyTimestamp = (timestamp) => {
  const date = new Date(timestamp)
  const YYYY = date.getFullYear()
  const MM = `${(date.getMonth() + 1 < 10 ? '0' : '') + (date.getMonth() + 1)}`
  const DD = `${(date.getDate() < 10 ? '0' : '') + date.getDate()}`
  const hh = `${(date.getHours() < 10 ? '0' : '') + date.getHours()}`
  const mm = `${(date.getMinutes() < 10 ? '0' : '') + date.getMinutes()}`
  const ss = `${(date.getSeconds() < 10 ? '0' : '') + date.getSeconds()}`
  return `${YYYY}-${MM}-${DD} ${hh}:${mm}:${ss}`
}
