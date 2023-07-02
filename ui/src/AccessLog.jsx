import { useState } from 'react'
import {
    makeStyles,
    Body1,
    Caption1,
    Button,
    shorthands,
    Card,
    CardFooter,
    CardHeader,
    CardPreview,
    Menu,
    MenuItem,
    MenuList,
    MenuPopover,
    MenuTrigger,
    SplitButton,
  } from "@fluentui/react-components"
import { CheckmarkRegular, DismissRegular, LocationRegular, RouterRegular, ClockRegular, NumberSymbolRegular } from "@fluentui/react-icons"
import Map, { Marker } from 'react-map-gl'
import 'mapbox-gl/dist/mapbox-gl.css'
import './AccessLog.css'
import api from './api'
import logo from './logo'
import dayjs from 'dayjs'

const useStyles = makeStyles({
  card: {
    ...shorthands.margin("auto"),
    width: "720px",
    maxWidth: "100%",
  },
  headerImage: {
    width: "32px",
    height: "32px",
  }
})

const checkConsent = (client, accessLog) => {
  for (let i = 0;i < client.consents.length;i++) {
    const consent = client.consents[i]
    if ((consent.host == '*' || consent.host == accessLog.host) && 
      (consent.ip == '*' || consent.ip == accessLog.ip_info.ip)) {
      if (consent.status == 'allowed') {
        return 'allowed'
      } else if (consent.status == 'blocked') {
        return 'blocked'
      }
      // else, continue to next consent
    }
  }
  return 'pending'
}

function AccessLog({client, accessLog}) {
  const styles = useStyles()
  const status = checkConsent(client, accessLog)

  const header = () => {
    const clientName = <b>{client.name}</b>
    return (
      <CardHeader
        image={
          <div className={styles.headerImage}>
            <img src={logo.browser(client)}/>
          </div>
        }
        header={
          <Body1>
            {clientName} wants to access <a target='_blank' href={`https://${accessLog.host}`}><b>{accessLog.host}</b></a>
          </Body1>
        }
        description={<Caption1>IP: {accessLog.ip_info.ip}</Caption1>}
      />)
  }

  const footer = () => {
    const actionAllow = () => {
      api.allowClient(client.id, accessLog.host)
    }
    const actionBlock = () => {
      api.blockClient(client.id, accessLog.host)
    }
    if (status === 'allowed') {
      return (
        <CardFooter>
          <Button onClick={actionAllow} icon={<CheckmarkRegular/>} appearance='primary'>Allowed</Button>
          <Button onClick={actionBlock} icon={<DismissRegular/>}>Block</Button>
        </CardFooter>
      )
    } else if (status === 'blocked') {
      return (
        <CardFooter>
          <Button onClick={actionAllow} icon={<CheckmarkRegular/>}>Allow</Button>
          <Button onClick={actionBlock} icon={<DismissRegular/>} appearance='primary'>Blocked</Button>
        </CardFooter>
      )
    }
    return (
      <CardFooter>
        <Button onClick={actionAllow} icon={<CheckmarkRegular/>}>Allow</Button>
        <Button onClick={actionBlock} icon={<DismissRegular/>}>Block</Button>
      </CardFooter>
    )
  }

  return (
    <div>
    <Card className={styles.card}>
      {header()}
      <CardPreview className='preview'>
        <div>
          <div className='map'>
            <Map
              mapboxAccessToken='pk.eyJ1IjoiaWxvYWh6IiwiYSI6ImNqd2dlZDM3MDFlb3E0OG84OGptZmx4YTYifQ.or3xgAAaDIzk3TNpl0rfWQ'
              mapLib={import('mapbox-gl')}
              initialViewState={{
                longitude: accessLog.ip_info.long - 1,
                latitude: accessLog.ip_info.lat,
                zoom: 6.5,
              }}
              style={{width: '100%', height: '100%'}}
              mapStyle="mapbox://styles/mapbox/streets-v9">
              <Marker longitude={accessLog.ip_info.long} latitude={accessLog.ip_info.lat} anchor='bottom'></Marker>
            </Map>
          </div>
          <div className='info'>
            <div className='line'>
              <LocationRegular className='icon' /> <span><a target='_blank' href={`https://maps.google.com/?q=${accessLog.ip_info.lat},${accessLog.ip_info.long}`}>{accessLog.ip_info.city}, {accessLog.ip_info.region}, {accessLog.ip_info.country}</a></span>
            </div>
            <div className='line'>
              <RouterRegular className='icon' /> <span><a target='_blank' href={`https://www.ip2location.com/as${accessLog.ip_info.asn}`}>{accessLog.ip_info.as}</a></span>
            </div>
            <div className='line'>
              <span className='line-logo'><img src={logo.os(client)} /></span>
              <span>{client.ua_info.os_name} {client.ua_info.os_version}</span>
              <span style={{width: '8px'}}></span>
              <span className='line-logo'><img src={logo.browser(client)} /></span>
              <span>{client.ua_info.browser_name} {client.ua_info.browser_version}</span>
            </div>
            <div className='line'>
              <ClockRegular className='icon' />
              <span>First seen {dayjs(accessLog.first_seen).format()}</span>
            </div>
            <div className='line'>
              <ClockRegular className='icon' />
              <span>Last seen {dayjs(accessLog.last_seen).format()}</span>
            </div>
            <div className='line'>
              <NumberSymbolRegular className='icon' />
              <span>Total visits {accessLog.count}</span>
            </div>
          </div>
        </div>
      </CardPreview>
      {footer()}
    </Card>
    </div>
  )
}

export default AccessLog
