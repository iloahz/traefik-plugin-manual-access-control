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
import './Client.css'
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

function Client({client}) {
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
            <a target='_blank' href={`https://www.ip2location.com/demo/${client.info.ip}`}><b>{client.info.ip}</b></a> wants to access <a target='_blank' href={`https://${client.url}`}><b>{client.url}</b></a>
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
                longitude: client.info.long - 1,
                latitude: client.info.lat,
                zoom: 6.5,
              }}
              style={{width: '100%', height: '100%'}}
              mapStyle="mapbox://styles/mapbox/streets-v9">
              <Marker longitude={client.info.long} latitude={client.info.lat} anchor='bottom'></Marker>
            </Map>
          </div>
          <div className='info'>
            <div className='line'>
              <LocationRegular className='icon' /> <span><a target='_blank' href={`https://maps.google.com/?q=${client.info.lat},${client.info.long}`}>{client.info.city}, {client.info.region}, {client.info.country}</a></span>
            </div>
            <div className='line'>
              <RouterRegular className='icon' /> <span><a target='_blank' href={`https://www.ip2location.com/as${client.info.asn}`}>{client.info.as}</a></span>
            </div>
            <div className='line'>
              <span className='line-logo'><img src={logo.os(client)} /></span>
              <span>{client.info.os_name} {client.info.os_version}</span>
              <span style={{width: '8px'}}></span>
              <span className='line-logo'><img src={logo.browser(client)} /></span>
              <span>{client.info.browser_name} {client.info.browser_version}</span>
            </div>
            <div className='line'>
              <ClockRegular className='icon' />
              <span>First seen {dayjs(client.stats.first_seen).format()}</span>
            </div>
            <div className='line'>
              <ClockRegular className='icon' />
              <span>Last seen {dayjs(client.stats.last_seen).format()}</span>
            </div>
            <div className='line'>
              <NumberSymbolRegular className='icon' />
              <span>Total visits {client.stats.count}</span>
            </div>
          </div>
        </div>
      </CardPreview>
      {footer()}
    </Card>
    </div>
  )
}

export default Client
