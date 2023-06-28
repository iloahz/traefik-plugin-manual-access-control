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

function AccessLog({client, accessLog}) {
  const styles = useStyles();

  const header = () => {
    return (
      <CardHeader
        image={
          <div className={styles.headerImage}>
            <img src={logo.browser(client)}/>
          </div>
        }
        header={
          <Body1>
            <a target='_blank' href={`https://www.ip2location.com/demo/${accessLog.ip_info.ip}`}><b>{accessLog.ip_info.ip}</b></a> wants to access <a target='_blank' href={`https://${accessLog.host}`}><b>{accessLog.host}</b></a>
          </Body1>
        }
        description={<Caption1>Client ID: {client.id}</Caption1>}
      />)
  }

  const footer = () => {
    if (client.status === 'allowed') {
      return (
        <CardFooter>
          <Button onClick={() => {api.allowClient(client.id)}} icon={<CheckmarkRegular/>} appearance='primary'>Allowed</Button>
          <Button onClick={() => {api.blockClient(client.id)}} icon={<DismissRegular/>}>Block</Button>
        </CardFooter>
      )
    } else if (client.status === 'blocked') {
      return (
        <CardFooter>
          <Button onClick={() => {api.allowClient(client.id)}} icon={<CheckmarkRegular/>}>Allow</Button>
          <Button onClick={() => {api.blockClient(client.id)}} icon={<DismissRegular/>} appearance='primary'>Blocked</Button>
        </CardFooter>
      )
    } else if (client.status === 'pending') {
      return (
        <CardFooter>
          <Button onClick={() => {api.allowClient(client.id)}} icon={<CheckmarkRegular/>}>Allow</Button>
          <Button onClick={() => {api.blockClient(client.id)}} icon={<DismissRegular/>}>Block</Button>
        </CardFooter>
      )
    }
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
