"use client"

import React, { useState } from "react"
import { BrutalistButton } from "@/components/examples/brutalist-button"
import { Card, CardContent, CardFooter, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { cn } from "@/lib/utils"
import './styles.css'

export default function BrutalistButtonDemo() {
  const [count, setCount] = useState(0)
  const [isLoading, setIsLoading] = useState(false)

  const handleIncrement = () => setCount(prev => prev + 1)
  const handleDecrement = () => setCount(prev => Math.max(0, prev - 1))
  const handleReset = () => setCount(0)

  const simulateLoading = () => {
    setIsLoading(true)
    setTimeout(() => setIsLoading(false), 2000)
  }

  return (
    <div className="container mx-auto py-10 px-4">
      <h1 className="text-4xl font-bold mb-8 font-mono">BrutalistButton Demo</h1>
      
      <section className="demo-section">
        <h2 className="demo-section-title">Button Variants</h2>
        <p className="demo-section-description">
          The BrutalistButton component comes with multiple variants: default, primary, secondary, success, warning, danger, ghost, and link.
          Each variant has distinct styling to communicate different purposes.
        </p>
        
        <div className="demo-grid">
          <Card>
            <CardHeader>
              <CardTitle>Default</CardTitle>
              <CardDescription>The standard button style</CardDescription>
            </CardHeader>
            <CardContent className="flex justify-center">
              <BrutalistButton>Default Button</BrutalistButton>
            </CardContent>
            <CardFooter className="text-sm text-muted-foreground">
              <code>{`<BrutalistButton>Default Button</BrutalistButton>`}</code>
            </CardFooter>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>Danger</CardTitle>
              <CardDescription>For destructive or high-risk actions</CardDescription>
            </CardHeader>
            <CardContent className="flex justify-center">
              <BrutalistButton variant="danger">Danger Button</BrutalistButton>
            </CardContent>
            <CardFooter className="text-sm text-muted-foreground">
              <code>{`<BrutalistButton variant="danger">Danger Button</BrutalistButton>`}</code>
            </CardFooter>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>Success</CardTitle>
              <CardDescription>For positive or confirming actions</CardDescription>
            </CardHeader>
            <CardContent className="flex justify-center">
              <BrutalistButton variant="success">Success Button</BrutalistButton>
            </CardContent>
            <CardFooter className="text-sm text-muted-foreground">
              <code>{`<BrutalistButton variant="success">Success Button</BrutalistButton>`}</code>
            </CardFooter>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>Warning</CardTitle>
              <CardDescription>For actions that need caution</CardDescription>
            </CardHeader>
            <CardContent className="flex justify-center">
              <BrutalistButton variant="warning">Warning Button</BrutalistButton>
            </CardContent>
            <CardFooter className="text-sm text-muted-foreground">
              <code>{`<BrutalistButton variant="warning">Warning Button</BrutalistButton>`}</code>
            </CardFooter>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>Ghost</CardTitle>
              <CardDescription>A subtle, transparent button style</CardDescription>
            </CardHeader>
            <CardContent className="flex justify-center">
              <BrutalistButton variant="ghost">Ghost Button</BrutalistButton>
            </CardContent>
            <CardFooter className="text-sm text-muted-foreground">
              <code>{`<BrutalistButton variant="ghost">Ghost Button</BrutalistButton>`}</code>
            </CardFooter>
          </Card>
        </div>
      </section>
      
      <section className="demo-section">
        <h2 className="demo-section-title">Button Sizes</h2>
        <p className="demo-section-description">
          BrutalistButton provides three size options: small, default (medium), and large.
          Choose the appropriate size based on the context and importance of the action.
        </p>
        
        <div className="demo-grid">
          <Card>
            <CardHeader>
              <CardTitle>Small</CardTitle>
              <CardDescription>Compact button for tight spaces</CardDescription>
            </CardHeader>
            <CardContent className="flex justify-center">
              <BrutalistButton size="sm">Small Button</BrutalistButton>
            </CardContent>
            <CardFooter className="text-sm text-muted-foreground">
              <code>{`<BrutalistButton size="sm">Small Button</BrutalistButton>`}</code>
            </CardFooter>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>Default</CardTitle>
              <CardDescription>Standard button size</CardDescription>
            </CardHeader>
            <CardContent className="flex justify-center">
              <BrutalistButton>Default Size</BrutalistButton>
            </CardContent>
            <CardFooter className="text-sm text-muted-foreground">
              <code>{`<BrutalistButton>Default Size</BrutalistButton>`}</code>
            </CardFooter>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>Large</CardTitle>
              <CardDescription>Prominent button for important actions</CardDescription>
            </CardHeader>
            <CardContent className="flex justify-center">
              <BrutalistButton size="lg">Large Button</BrutalistButton>
            </CardContent>
            <CardFooter className="text-sm text-muted-foreground">
              <code>{`<BrutalistButton size="lg">Large Button</BrutalistButton>`}</code>
            </CardFooter>
          </Card>
        </div>
      </section>
      
      <section className="demo-section">
        <h2 className="demo-section-title">Interactive Examples</h2>
        <p className="demo-section-description">
          BrutalistButton components are fully interactive and can be used with state management.
          Here are some examples of buttons with dynamic behaviors.
        </p>
        
        <div className="demo-grid">
          <Card>
            <CardHeader>
              <CardTitle>Counter Example</CardTitle>
              <CardDescription>Click the buttons to change the count</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="interactive-example">
                <div className="counter-display">{count}</div>
                <div className="button-container">
                  <BrutalistButton variant="success" onClick={handleIncrement}>Increment</BrutalistButton>
                  <BrutalistButton variant="danger" onClick={handleDecrement}>Decrement</BrutalistButton>
                  <BrutalistButton variant="ghost" onClick={handleReset}>Reset</BrutalistButton>
                </div>
              </div>
            </CardContent>
            <CardFooter className="text-sm text-muted-foreground">
              Buttons with onClick handlers to update state
            </CardFooter>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>Loading State</CardTitle>
              <CardDescription>Button with loading indicator</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="interactive-example">
                <div className="loading-button-container">
                  <BrutalistButton 
                    variant="warning" 
                    onClick={simulateLoading}
                    disabled={isLoading}
                  >
                    {isLoading ? (
                      <>
                        <span className="loading-indicator"></span>
                        Loading...
                      </>
                    ) : (
                      'Click to Load'
                    )}
                  </BrutalistButton>
                  <p className="text-sm text-center">
                    {isLoading ? 'Processing...' : 'Ready'}
                  </p>
                </div>
              </div>
            </CardContent>
            <CardFooter className="text-sm text-muted-foreground">
              Demonstrates disabled state and loading indicators
            </CardFooter>
          </Card>
        </div>
      </section>
      
      <section className="demo-section">
        <h2 className="demo-section-title">Shadow & Transform Effects</h2>
        <p className="demo-section-description">
          BrutalistButton supports various shadow and transform effects to create a dynamic, interactive feel.
        </p>
        
        <div className="demo-grid">
          <Card>
            <CardHeader>
              <CardTitle>Shadow Effects</CardTitle>
              <CardDescription>Buttons with different shadow properties</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="button-container">
                <div className="button-example">
                  <BrutalistButton shadow="none">No Shadow</BrutalistButton>
                  <span className="button-label">shadow="none"</span>
                </div>
                <div className="button-example">
                  <BrutalistButton shadow="sm">Small Shadow</BrutalistButton>
                  <span className="button-label">shadow="sm"</span>
                </div>
                <div className="button-example">
                  <BrutalistButton>Default Shadow</BrutalistButton>
                  <span className="button-label">default</span>
                </div>
                <div className="button-example">
                  <BrutalistButton shadow="lg">Large Shadow</BrutalistButton>
                  <span className="button-label">shadow="lg"</span>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>Transform Effects</CardTitle>
              <CardDescription>Buttons with different transform properties</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="button-container">
                <div className="button-example">
                  <BrutalistButton transform="none">No Transform</BrutalistButton>
                  <span className="button-label">transform="none"</span>
                </div>
                <div className="button-example">
                  <BrutalistButton>Default Transform</BrutalistButton>
                  <span className="button-label">default</span>
                </div>
                <div className="button-example">
                  <BrutalistButton transform="lg">Large Transform</BrutalistButton>
                  <span className="button-label">transform="lg"</span>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>Combined Features</CardTitle>
              <CardDescription>Buttons with multiple customizations</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="button-container">
                <div className="button-example">
                  <BrutalistButton 
                    variant="success"
                    size="lg"
                    shadow="lg"
                    transform="lg"
                  >
                    Custom Button
                  </BrutalistButton>
                  <span className="button-label">Multiple props</span>
                </div>
                <div className="button-example">
                  <BrutalistButton 
                    variant="ghost"
                    size="sm"
                    shadow="none"
                    transform="none"
                  >
                    Subtle Button
                  </BrutalistButton>
                  <span className="button-label">Multiple props</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </section>
      
      <section className="demo-section usage-example">
        <h2 className="demo-section-title">Usage Example</h2>
        <p className="demo-section-description">
          Here's how to use the BrutalistButton component in your code:
        </p>
        
        <div className="code-block">
          <code>{`import { BrutalistButton } from '@/components/examples/brutalist-button';

export default function MyComponent() {
  const handleClick = () => {
    console.log('Button clicked!');
  };

  return (
    <div>
      {/* Basic button */}
      <BrutalistButton>Click Me</BrutalistButton>
      
      {/* With variant */}
      <BrutalistButton variant="success">Success Action</BrutalistButton>
      
      {/* With size */}
      <BrutalistButton size="lg">Large Button</BrutalistButton>
      
      {/* With event handler */}
      <BrutalistButton onClick={handleClick}>Handle Click</BrutalistButton>
      
      {/* With multiple customizations */}
      <BrutalistButton
        variant="warning"
        size="sm"
        shadow="lg"
        transform="lg"
        onClick={() => alert('Warning action!')}
      >
        Custom Warning
      </BrutalistButton>
    </div>
  );
}`}</code>
        </div>
      </section>
    </div>
  )
} 